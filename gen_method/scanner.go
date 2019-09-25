package gen_method

import (
	"go/parser"
	"go/types"

	"golang.org/x/tools/go/loader"
)

func NewScanner(packagePath string) *Scanner {
	ldr := loader.Config{}
	ldr.AllowErrors = true

	ldr.ParserMode = parser.ParseComments
	ldr.Import(packagePath)

	prog, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	return &Scanner{
		packagePath: packagePath,
		prog:        prog,
	}
}

type Scanner struct {
	packagePath string
	prog        *loader.Program
}

func (s *Scanner) Output(modelName string, ignoreTable bool) {
	for pkg, packageInfo := range s.prog.AllPackages {
		if packageInfo.Pkg.Path() == s.packagePath {
			for ident, def := range packageInfo.Defs {
				if ident.Name == modelName {
					if typeName, ok := def.(*types.TypeName); ok {
						if typeStruct, ok := typeName.Type().Underlying().(*types.Struct); ok {
							m := Model{
								Pkg:            pkg,
								Name:           modelName,
								UniqueIndex:    make(map[string][]Field),
								NormalIndex:    make(map[string][]Field),
								FuncMapContent: make(map[string]string),
							}

							m.collectInfoFromStructType(typeStruct)
							m.Output(packageInfo.Pkg.Name(), ignoreTable)
						}
					}
				}
			}
		}
	}
}
