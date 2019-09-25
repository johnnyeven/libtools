package gen

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen/loaderx"
	"github.com/johnnyeven/libtools/courier"
)

var (
	courierPkgImportPath = "github.com/johnnyeven/libtools/courier"
	routerTypeString     = reflectTypeString(reflect.TypeOf(new(courier.Router)))
)

func isRouterType(tpe types.Type) bool {
	return strings.HasSuffix(tpe.String(), routerTypeString)
}

func NewRouterScanner(program *loader.Program) *RouterScanner {
	routerScanner := &RouterScanner{
		program: program,
		routers: map[*types.Var]*Router{},
	}

	routerScanner.init()

	return routerScanner
}

type RouterScanner struct {
	program *loader.Program
	routers map[*types.Var]*Router
}

func (scanner *RouterScanner) Router(typeName *types.Var) *Router {
	return scanner.routers[typeName]
}

type OperatorTypeNames []*types.TypeName

func (ops OperatorTypeNames) Append(tpe types.Type) OperatorTypeNames {
	switch tpe.(type) {
	case *types.Named:
		return append(ops, tpe.(*types.Named).Obj())
	case *types.Pointer:
		return append(ops, tpe.(*types.Pointer).Elem().(*types.Named).Obj())
	}
	return ops
}

func FromArgs(pkgInfo *loader.PackageInfo, args ...ast.Expr) OperatorTypeNames {
	opTypeNames := OperatorTypeNames{}
	for _, arg := range args {
		opTypeNames = opTypeNames.Append(pkgInfo.TypeOf(arg))
	}
	return opTypeNames
}

func (scanner *RouterScanner) init() {
	for _, pkgInfo := range scanner.program.AllPackages {
		for ident, obj := range pkgInfo.Defs {
			if typeVar, ok := obj.(*types.Var); ok {
				if typeVar != nil && !strings.HasSuffix(typeVar.Pkg().Path(), courierPkgImportPath) {
					if isRouterType(typeVar.Type()) {
						router := NewRouter()

						ast.Inspect(ident.Obj.Decl.(ast.Node), func(node ast.Node) bool {
							switch node.(type) {
							case *ast.CallExpr:
								callExpr := node.(*ast.CallExpr)
								router.AppendOperators(FromArgs(pkgInfo, callExpr.Args...)...)
								return false
							}
							return true
						})

						scanner.routers[typeVar] = router
					}
				}
			}
		}
	}

	for _, pkgInfo := range scanner.program.AllPackages {
		for selectExpr, selection := range pkgInfo.Selections {
			if selection.Obj() != nil {
				if typeFunc, ok := selection.Obj().(*types.Func); ok {
					recv := typeFunc.Type().(*types.Signature).Recv()
					if recv != nil && isRouterType(recv.Type()) {
						for typeVar, router := range scanner.routers {
							switch selectExpr.Sel.Name {
							case "Register":
								if typeVar == pkgInfo.ObjectOf(IdentOfCallExprSelectExpr(selectExpr)) {
									file := loaderx.FileOf(selectExpr, pkgInfo.Files...)
									ast.Inspect(file, func(node ast.Node) bool {
										switch node.(type) {
										case *ast.CallExpr:
											callExpr := node.(*ast.CallExpr)
											if callExpr.Fun == selectExpr {
												routerIdent := callExpr.Args[0]
												switch routerIdent.(type) {
												case *ast.Ident:
													argTypeVar := pkgInfo.ObjectOf(routerIdent.(*ast.Ident)).(*types.Var)
													if r, ok := scanner.routers[argTypeVar]; ok {
														router.Register(r)
													}
												case *ast.SelectorExpr:
													argTypeVar := pkgInfo.ObjectOf(routerIdent.(*ast.SelectorExpr).Sel).(*types.Var)
													if r, ok := scanner.routers[argTypeVar]; ok {
														router.Register(r)
													}
												case *ast.CallExpr:
													callExprForRegister := routerIdent.(*ast.CallExpr)
													router.With(FromArgs(pkgInfo, callExprForRegister.Args...)...)
												}
												return false
											}
										}
										return true
									})
								}
							}
						}
					}
				}
			}
		}
	}
}

func IdentOfCallExprSelectExpr(selectExpr *ast.SelectorExpr) *ast.Ident {
	switch selectExpr.X.(type) {
	case *ast.Ident:
		return selectExpr.X.(*ast.Ident)
	case *ast.SelectorExpr:
		return selectExpr.X.(*ast.SelectorExpr).Sel
	}
	return nil
}

func NewRouter(operators ...*types.TypeName) *Router {
	return &Router{
		operators: operators,
	}
}

type Router struct {
	parent    *Router
	operators []*types.TypeName
	children  map[*Router]bool
}

func (router *Router) AppendOperators(operators ...*types.TypeName) {
	router.operators = append(router.operators, operators...)
}

func (router *Router) With(operators ...*types.TypeName) {
	router.Register(NewRouter(operators...))
}

func (router *Router) Register(r *Router) {
	if router.children == nil {
		router.children = map[*Router]bool{}
	}
	r.parent = router
	router.children[r] = true
}

func (router *Router) Route(program *loader.Program) *Route {
	parent := router.parent
	operators := router.operators

	for parent != nil {
		operators = append(parent.operators, operators...)
		parent = parent.parent
	}

	route := Route{
		last:      router.children == nil,
		operators: operators,
	}

	route.SetMethod(program)
	route.SetPath(program)

	return &route
}

func (router *Router) Routes(program *loader.Program) (routes []*Route) {
	for child := range router.children {
		route := child.Route(program)
		if route.last {
			routes = append(routes, route)
		}
		if child.children != nil {
			routes = append(routes, child.Routes(program)...)
		}
	}
	return routes
}

type Route struct {
	Method    string
	Path      string
	last      bool
	operators []*types.TypeName
}

func (route *Route) SetPath(program *loader.Program) {
	p := "/"
	for _, operator := range route.operators {
		typeFunc := loaderx.MethodOf(operator.Type().(*types.Named), "Path")

		loaderx.ForEachFuncResult(program, typeFunc, func(resultTypeAndValues ...types.TypeAndValue) {
			if resultTypeAndValues[0].IsValue() {
				p += getConstVal(resultTypeAndValues[0].Value).(string)
			}
		})
	}
	route.Path = httprouter.CleanPath(p)
}

func (route *Route) SetMethod(program *loader.Program) {
	if len(route.operators) > 0 {
		operator := route.operators[len(route.operators)-1]
		typeFunc := loaderx.MethodOf(operator.Type().(*types.Named), "Method")

		loaderx.ForEachFuncResult(program, typeFunc, func(resultTypeAndValues ...types.TypeAndValue) {
			if resultTypeAndValues[0].IsValue() {
				route.Method = getConstVal(resultTypeAndValues[0].Value).(string)
			}
		})
	}
}
