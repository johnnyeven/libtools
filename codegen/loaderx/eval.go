package loaderx

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
)

func StringifyAst(fset *token.FileSet, node ast.Node) string {
	buf := bytes.Buffer{}
	if err := format.Node(&buf, fset, node); err != nil {
		panic(err)
	}
	return buf.String()
}

func MustEvalExpr(fileSet *token.FileSet, pkg *types.Package, expr ast.Expr) types.TypeAndValue {
	if expr == nil {
		return types.TypeAndValue{
			Type: types.Typ[types.UntypedNil],
		}
	}

	if ident, ok := expr.(*ast.Ident); ok && ident.Name == "nil" {
		return types.TypeAndValue{
			Type: types.Typ[types.UntypedNil],
		}
	}

	code := StringifyAst(fileSet, expr)
	tv, err := types.Eval(fileSet, pkg, expr.Pos(), code)
	if err != nil {
		panic(err)
	}
	return tv
}
