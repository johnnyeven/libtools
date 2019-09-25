package loaderx

import (
	"go/types"

	"golang.org/x/tools/go/loader"
)

func NewQuery(program *loader.Program, pkgImportPath string) *Query {
	pkgInfo := program.Package(pkgImportPath)

	return &Query{
		pkgInfo: pkgInfo,
	}
}

type Query struct {
	pkgInfo *loader.PackageInfo
}

func (q *Query) Const(name string) *types.Const {
	for ident, def := range q.pkgInfo.Defs {
		if typeConst, ok := def.(*types.Const); ok {
			if ident.Name == name {
				return typeConst
			}
		}
	}
	return nil
}

func (q *Query) TypeName(name string) *types.TypeName {
	for ident, def := range q.pkgInfo.Defs {
		if typeName, ok := def.(*types.TypeName); ok {
			if ident.Name == name {
				return typeName
			}
		}
	}
	return nil
}

func (q *Query) Var(name string) *types.Var {
	for ident, def := range q.pkgInfo.Defs {
		if typeVar, ok := def.(*types.Var); ok {
			if ident.Name == name {
				return typeVar
			}
		}
	}
	return nil
}

func (q *Query) Func(name string) *types.Func {
	for ident, def := range q.pkgInfo.Defs {
		if typeFunc, ok := def.(*types.Func); ok {
			if ident.Name == name {
				return typeFunc
			}
		}
	}
	return nil
}
