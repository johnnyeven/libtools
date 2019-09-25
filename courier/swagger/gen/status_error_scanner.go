package gen

import (
	"go/ast"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen/loaderx"
	"github.com/johnnyeven/libtools/courier/status_error"
	"github.com/johnnyeven/libtools/courier/status_error/gen"
)

var (
	statusErrorTypeString     = reflectTypeString(reflect.TypeOf(new(status_error.StatusError)))
	statusErrorCodeTypeString = reflectTypeString(reflect.TypeOf(new(status_error.StatusErrorCode)))
)

func NewStatusErrorScanner(program *loader.Program) *StatusErrorScanner {
	statusErrorScanner := &StatusErrorScanner{
		program:      program,
		statusErrors: map[*types.Const]status_error.StatusError{},
		errorsUsed:   map[*types.Func]status_error.StatusErrorCodeMap{},
	}

	statusErrorScanner.init()

	return statusErrorScanner
}

type StatusErrorScanner struct {
	program      *loader.Program
	statusErrors map[*types.Const]status_error.StatusError
	errorsUsed   map[*types.Func]status_error.StatusErrorCodeMap
}

func (scanner *StatusErrorScanner) StatusErrorsInFunc(typeFunc *types.Func) status_error.StatusErrorCodeMap {
	if typeFunc == nil {
		return nil
	}

	if statusErrorCodeMap, ok := scanner.errorsUsed[typeFunc]; ok {
		return statusErrorCodeMap
	}

	// force initial to fix recursion
	scanner.errorsUsed[typeFunc] = status_error.StatusErrorCodeMap{}

	pkgInfo := scanner.program.Package(typeFunc.Pkg().Path())
	funcDecl := loaderx.FuncDeclOfTypeFunc(pkgInfo, typeFunc)

	if funcDecl != nil {
		ast.Inspect(funcDecl, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.CallExpr:
				identList := loaderx.GetIdentChainCallOfCallFun(node.(*ast.CallExpr).Fun)
				if len(identList) > 0 {
					callIdent := identList[len(identList)-1]
					obj := pkgInfo.ObjectOf(callIdent)

					if nextTypeFunc, ok := obj.(*types.Func); ok && nextTypeFunc != typeFunc && nextTypeFunc.Pkg() != nil {
						statusErrorCodeMap := scanner.StatusErrorsInFunc(nextTypeFunc)
						scanner.mayMergeStateError(typeFunc, statusErrorCodeMap)
					}
				}
			case *ast.Ident:
				scanner.mayAddStateErrorByObject(typeFunc, pkgInfo.ObjectOf(node.(*ast.Ident)))
			}
			return true
		})

		doc := loaderx.StringifyCommentGroup(funcDecl.Doc)
		scanner.mayMergeStateError(typeFunc, pickStatusErrorsFromDoc(doc))
	}

	return scanner.errorsUsed[typeFunc]
}

func (scanner *StatusErrorScanner) mayAddStateErrorByObject(typeFunc *types.Func, obj types.Object) {
	if obj == nil {
		return
	}
	if typeConst, ok := obj.(*types.Const); ok {
		scanner.mayAddStateError(typeFunc, typeConst)
	}
}

func (scanner *StatusErrorScanner) mayMergeStateError(typeFunc *types.Func, statusErrorCodeMap status_error.StatusErrorCodeMap) {
	if scanner.errorsUsed[typeFunc] == nil {
		scanner.errorsUsed[typeFunc] = status_error.StatusErrorCodeMap{}
	}
	scanner.errorsUsed[typeFunc].Merge(statusErrorCodeMap)
}

func (scanner *StatusErrorScanner) mayAddStateError(typeFunc *types.Func, typeConst *types.Const) {
	if statusError, ok := scanner.statusErrors[typeConst]; ok {
		if scanner.errorsUsed[typeFunc] == nil {
			scanner.errorsUsed[typeFunc] = status_error.StatusErrorCodeMap{}
		}
		scanner.errorsUsed[typeFunc][int64(statusError.Code)] = statusError
	}
}

func (scanner *StatusErrorScanner) init() {
	for _, pkgInfo := range scanner.program.AllPackages {
		for ident, obj := range pkgInfo.Defs {
			if constObj, ok := obj.(*types.Const); ok {
				if constObj.Type().String() == statusErrorCodeTypeString {
					key := constObj.Name()
					if key == "_" {
						continue
					}

					doc := loaderx.CommentsOf(scanner.program.Fset, ident, pkgInfo.Files...)
					code, _ := strconv.ParseInt(constObj.Val().String(), 10, 64)
					msg, desc, canBeErrTalk := gen.ParseStatusErrorDesc(doc)

					scanner.statusErrors[constObj] = status_error.StatusError{
						Key:            key,
						Code:           code,
						Msg:            msg,
						Desc:           desc,
						CanBeErrorTalk: canBeErrTalk,
					}
				}
			}
		}
	}
}

func pickStatusErrorsFromDoc(doc string) status_error.StatusErrorCodeMap {
	statusErrorCodeMap := status_error.StatusErrorCodeMap{}

	lines := strings.Split(doc, "\n")

	for _, line := range lines {
		if line != "" {
			if statusErr := status_error.ParseString(line); statusErr != nil {
				statusErrorCodeMap[statusErr.Code] = *statusErr
			}
		}
	}
	return statusErrorCodeMap
}
