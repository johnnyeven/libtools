package loaderx

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/loader"
)

func FuncDeclOfTypeFunc(pkgInfo *loader.PackageInfo, typeFunc *types.Func) *ast.FuncDecl {
	for ident, def := range pkgInfo.Defs {
		if typeFuncDef, ok := def.(*types.Func); ok {
			if typeFuncDef == typeFunc {
				return FuncDeclOf(ident, FileOf(ident, pkgInfo.Files...))
			}
		}
	}
	return nil
}

func FuncDeclOf(ident *ast.Ident, file *ast.File) (funcDecl *ast.FuncDecl) {
	ast.Inspect(file, func(node ast.Node) bool {
		if decl, ok := node.(*ast.FuncDecl); ok {
			if decl.Name == ident {
				funcDecl = decl
				return false
			}
		}
		return true
	})
	return
}

func GetIdentChainCallOfCallFun(expr ast.Expr) (list []*ast.Ident) {
	switch expr.(type) {
	case *ast.SelectorExpr:
		selectorExpr := expr.(*ast.SelectorExpr)
		list = append(list, GetIdentChainCallOfCallFun(selectorExpr.X)...)
		list = append(list, selectorExpr.Sel)
	case *ast.Ident:
		list = append(list, expr.(*ast.Ident))
	}
	return
}

func ForEachFuncResult(program *loader.Program, typeFunc *types.Func, walker func(resultTypeAndValues ...types.TypeAndValue)) {
	if typeFunc == nil {
		return
	}

	pkgInfo := program.Package(typeFunc.Pkg().Path())
	funcDecl := FuncDeclOfTypeFunc(pkgInfo, typeFunc)

	if funcDecl == nil {
		// todo find way to location interface
		return
	}

	signature := typeFunc.Type().(*types.Signature)
	results := signature.Results()
	resultLength := results.Len()

	returnStmtList := make([]*ast.ReturnStmt, 0)

	// collect all return stmt
	ast.Inspect(funcDecl, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.FuncLit:
			// skip func inline declaration
			return false
		case *ast.ReturnStmt:
			returnStmtList = append(returnStmtList, node.(*ast.ReturnStmt))
		}
		return true
	})

	for _, returnStmt := range returnStmtList {
		namedResults := make([]ast.Expr, resultLength)
		typeAndValues := make([]types.TypeAndValue, resultLength)
		skip := false

		ast.Inspect(funcDecl, func(node ast.Node) bool {
			switch node.(type) {
			// skip for `switch-case` if return no in
			case *ast.CaseClause:
				caseBody := node.(*ast.CaseClause)
				return !hasReturn(node) || caseBody.Pos() <= returnStmt.Pos() && returnStmt.Pos() < caseBody.End()
			case *ast.IfStmt:
				ifBody := node.(*ast.IfStmt).Body
				return !hasReturn(node) || ifBody.Pos() <= returnStmt.Pos() && returnStmt.Pos() < ifBody.End()
			case *ast.ReturnStmt:
				currentReturnStmt := node.(*ast.ReturnStmt)
				if !skip && currentReturnStmt == returnStmt {
					if currentReturnStmt.Results == nil {
						for i := 0; i < resultLength; i++ {
							typeAndValues[i] = MustEvalExpr(program.Fset, pkgInfo.Pkg, namedResults[i])
							typeAndValues[i] = patchTypeAndValue(results.At(i).Type(), typeAndValues[i])
						}
						walker(typeAndValues...)
					} else {
						if len(currentReturnStmt.Results) < resultLength {
							if callExpr, ok := currentReturnStmt.Results[0].(*ast.CallExpr); ok {
								identList := GetIdentChainCallOfCallFun(callExpr.Fun)
								ForEachFuncResult(program, pkgInfo.ObjectOf(identList[len(identList)-1]).(*types.Func), walker)
							}
						} else {
							for i := 0; i < resultLength; i++ {
								typeAndValues[i] = MustEvalExpr(program.Fset, pkgInfo.Pkg, currentReturnStmt.Results[i])
								typeAndValues[i] = patchTypeAndValue(results.At(i).Type(), typeAndValues[i])
							}
							walker(typeAndValues...)
						}
					}
				}
			case *ast.AssignStmt:
				// only scan before return
				if node != nil && node.Pos() >= returnStmt.Pos() {
					return false
				}

				assignStmt := node.(*ast.AssignStmt)

				if len(assignStmt.Rhs) == 1 {
					if _, ok := assignStmt.Rhs[0].(*ast.FuncLit); ok {
						return false
					}
				}

				if len(assignStmt.Lhs) == len(assignStmt.Rhs) {
					for i, expr := range assignStmt.Lhs {
						if ident, ok := expr.(*ast.Ident); ok {
							for resultIndex := 0; resultIndex < resultLength; resultIndex++ {
								if pkgInfo.ObjectOf(ident) == results.At(resultIndex) {
									namedResults[resultIndex] = assignStmt.Rhs[i]
								}
							}
						}
					}
				} else {
					if callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
						allNamedResult := false
						for resultIndex := 0; resultIndex < resultLength; resultIndex++ {
							switch lhs := assignStmt.Lhs[0].(type) {
							case *ast.Ident:
								if pkgInfo.ObjectOf(lhs) == results.At(resultIndex) {
									allNamedResult = true
								}
							case *ast.SelectorExpr:
								if pkgInfo.ObjectOf(lhs.Sel) == results.At(resultIndex) {
									allNamedResult = true
								}
							}
						}
						if allNamedResult {
							identList := GetIdentChainCallOfCallFun(callExpr.Fun)
							skip = true
							ForEachFuncResult(program, pkgInfo.ObjectOf(identList[len(identList)-1]).(*types.Func), walker)
						}
					}
				}
			}
			return true
		})
	}
}

func hasReturn(node ast.Node) (ok bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.ReturnStmt:
			ok = true
		}
		return true
	})
	return
}

func patchTypeAndValue(tpe types.Type, typeAndValue types.TypeAndValue) types.TypeAndValue {
	if typeAndValue.IsValue() && typeAndValue.Value == nil {
		return types.TypeAndValue{
			Type: typeAndValue.Type,
		}
	}
	_, isInterface := tpe.(*types.Interface)
	if !isInterface && typeAndValue.Type == types.Typ[types.UntypedNil] {
		return types.TypeAndValue{
			Type: tpe,
		}
	}
	return typeAndValue
}

func MethodOf(named *types.Named, funcName string) (typeFunc *types.Func) {
	for i := 0; i < named.NumMethods(); i++ {
		method := named.Method(i)
		if method.Name() == funcName {
			return method
		}
	}

	if structType, ok := named.Underlying().(*types.Struct); ok {
		for i := 0; i < structType.NumFields(); i++ {
			field := structType.Field(i)
			if field.Anonymous() {
				typeFuncAnonymous := MethodOf(IndirectType(field.Type()).(*types.Named), funcName)
				if typeFunc == nil {
					typeFunc = typeFuncAnonymous
				} else if typeFuncAnonymous != nil {
					typeFunc = nil
				}
			}
		}
	}
	return
}

func IndirectType(tpe types.Type) types.Type {
	switch tpe.(type) {
	case *types.Pointer:
		return IndirectType(tpe.(*types.Pointer).Elem())
	default:
		return tpe
	}
}
