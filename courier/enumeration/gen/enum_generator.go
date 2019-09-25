package gen

import (
	"go/build"
	"go/parser"
	"go/types"
	"path"
	"path/filepath"

	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/codegen/loaderx"
	"github.com/johnnyeven/libtools/courier/swagger/gen"
	"github.com/johnnyeven/libtools/godash"
)

type EnumGenerator struct {
	Filters       []string
	pkgImportPath string
	program       *loader.Program
	enumScanner   *gen.EnumScanner
}

var _ interface {
	codegen.Generator
} = (*EnumGenerator)(nil)

func (g *EnumGenerator) Load(cwd string) {
	ldr := loader.Config{
		AllowErrors: true,
		ParserMode:  parser.ParseComments,
	}

	pkgImportPath := codegen.GetPackageImportPath(cwd)
	ldr.Import(pkgImportPath)

	p, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	g.program = p
	g.pkgImportPath = pkgImportPath
	g.enumScanner = gen.NewEnumScanner(p)
}

func (g *EnumGenerator) Pick() {
	for pkg, pkgInfo := range g.program.AllPackages {
		if pkg.Path() != g.pkgImportPath {
			continue
		}
		for ident, obj := range pkgInfo.Defs {
			if typeName, ok := obj.(*types.TypeName); ok {
				doc := loaderx.CommentsOf(g.program.Fset, ident, pkgInfo.Files...)
				doc, hasEnum := gen.ParseEnum(doc)

				if hasEnum {
					if len(g.Filters) > 0 {
						if godash.StringIncludes(g.Filters, typeName.Name()) {
							g.enumScanner.Enum(typeName)
						}
					} else {
						g.enumScanner.Enum(typeName)
					}
				}
			}
		}
	}
}

func (g *EnumGenerator) Output(cwd string) codegen.Outputs {
	outputs := codegen.Outputs{}
	for typeName, e := range g.enumScanner.Enums {
		p, _ := build.Import(typeName.Pkg().Path(), "", build.FindOnly)
		dir, _ := filepath.Rel(cwd, p.Dir)

		enum := NewEnum(typeName.Pkg().Path(), typeName.Pkg().Name(), typeName.Name(), e, g.enumScanner.HasOffset(typeName))

		outputs.Add(codegen.GeneratedSuffix(path.Join(dir, codegen.ToLowerSnakeCase(typeName.Name())+".go")), enum.String())
	}
	return outputs
}
