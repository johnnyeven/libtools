package gen

import (
	"fmt"
	"go/parser"
	"go/types"
	"regexp"
	"strings"

	"github.com/morlay/oas"
	"golang.org/x/tools/go/loader"

	"encoding/json"

	"github.com/johnnyeven/libtools/codegen"
)

type SwaggerGenerator struct {
	RootRouterName  string
	Next            bool
	pkgImportPath   string
	program         *loader.Program
	openapi         *oas.OpenAPI
	routerScanner   *RouterScanner
	operatorScanner *OperatorScanner
}

func (g *SwaggerGenerator) Load(cwd string) {
	ldr := loader.Config{
		ParserMode: parser.ParseComments,
	}

	pkgImportPath := codegen.GetPackageImportPath(cwd)
	ldr.Import(pkgImportPath)

	p, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	g.program = p
	g.pkgImportPath = pkgImportPath
	g.openapi = oas.NewOpenAPI()
	g.operatorScanner = NewOperatorScanner(p)
	g.routerScanner = NewRouterScanner(p)
}

func (g *SwaggerGenerator) Pick() {
	for _, pkgInfo := range g.program.AllPackages {
		for _, def := range pkgInfo.Defs {
			if typeVar, ok := def.(*types.Var); ok {
				if typeVar.Name() == g.RootRouterName {
					router := g.routerScanner.Router(typeVar)
					routes := router.Routes(g.program)
					operationIDs := map[string]*Route{}
					for _, route := range routes {
						operation := g.getOperationByOperatorTypes(route.Method, route.operators...)
						if _, exists := operationIDs[operation.OperationId]; exists {
							panic(fmt.Errorf("operationID %s should be unique", operation.OperationId))
						}
						operationIDs[operation.OperationId] = route
						g.openapi.AddOperation(oas.HttpMethod(strings.ToLower(route.Method)), g.patchPath(route.Path, operation), operation)
					}
					g.operatorScanner.BindSchemas(g.openapi)
					return
				}
			}
		}
	}
}

var RxHttpRouterPath = regexp.MustCompile("/:([^/]+)")

func (g *SwaggerGenerator) patchPath(swaggerPath string, operation *oas.Operation) string {
	return RxHttpRouterPath.ReplaceAllStringFunc(swaggerPath, func(str string) string {
		name := RxHttpRouterPath.FindAllStringSubmatch(str, -1)[0][1]

		var isParameterDefined = false

		for _, parameter := range operation.Parameters {
			if parameter.In == "path" && parameter.Name == name {
				isParameterDefined = true
			}
		}

		if isParameterDefined {
			return "/{" + name + "}"
		}

		return "/0"
	})
}

func (g *SwaggerGenerator) getOperationByOperatorTypes(method string, operatorTypes ...*types.TypeName) *oas.Operation {
	operators := make([]Operator, 0)

	for _, operatorType := range operatorTypes {
		operators = append(operators, *g.operatorScanner.Operator(operatorType))
	}

	return ConcatToOperation(method, operators...)
}

func (g *SwaggerGenerator) Output(cwd string) codegen.Outputs {
	bytes, _ := json.MarshalIndent(g.openapi, "", "  ")

	return codegen.Outputs{
		"swagger.json": string(bytes),
	}
}
