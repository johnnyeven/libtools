package gen_from_old

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/types"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/courier/status_error"
)

type StatusErrorGenerator struct {
	DryRun        bool
	pkgImportPath string
	program       *loader.Program
	statusErrors  map[int64]status_error.StatusError
}

func (g *StatusErrorGenerator) Load(cwd string) {
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
	g.statusErrors = map[int64]status_error.StatusError{}
}

func (g *StatusErrorGenerator) Pick() {
	statusErrorType := reflect.TypeOf(status_error.StatusError{})
	statusErrorTypeFullName := fmt.Sprintf("%s.%s", statusErrorType.PkgPath(), statusErrorType.Name())

	for pkg, pkgInfo := range g.program.AllPackages {
		if pkg.Path() != g.pkgImportPath {
			continue
		}
		for ident, obj := range pkgInfo.Defs {
			if varObj, ok := obj.(*types.Var); ok {
				if strings.HasSuffix(varObj.Type().String(), statusErrorTypeFullName) {
					key := varObj.Name()
					statusErr := status_error.StatusError{}
					statusErr.Key = key

					if valueSpec, ok := ident.Obj.Decl.(*ast.ValueSpec); ok {
						ast.Inspect(valueSpec, func(node ast.Node) bool {
							if keyValueExpr, ok := node.(*ast.KeyValueExpr); ok {
								switch keyValueExpr.Key.(*ast.Ident).Name {
								case "Code":
									if basicLit, ok := keyValueExpr.Value.(*ast.BasicLit); ok {
										statusErr.Code, _ = strconv.ParseInt(basicLit.Value, 10, 64)
									}
								case "Msg":
									if basicLit, ok := keyValueExpr.Value.(*ast.BasicLit); ok {
										statusErr.Msg, _ = strconv.Unquote(basicLit.Value)
									}
								case "Desc":
									if basicLit, ok := keyValueExpr.Value.(*ast.BasicLit); ok {
										statusErr.Desc, _ = strconv.Unquote(basicLit.Value)
									}
								case "CanBeErrorTalk":
									if ident, ok := keyValueExpr.Value.(*ast.Ident); ok {
										statusErr.CanBeErrorTalk = ident.Name == "true"
									}
								}
							}
							return true
						})
					}

					if s, ok := g.statusErrors[statusErr.Code]; ok {
						panic(fmt.Errorf("%d already used in %s", statusErr.Code, s.Error()))
					}
					g.statusErrors[statusErr.Code] = statusErr
				}
			}
		}
	}
}

func (g *StatusErrorGenerator) Output(cwd string) codegen.Outputs {
	outputs := codegen.Outputs{}
	codes := make([]int, 0)
	for code := range g.statusErrors {
		codes = append(codes, int(code))
	}
	sort.Ints(codes)

	statusErrorGroups := make(map[int][]status_error.StatusError)

	for _, code := range codes {
		statueErr := g.statusErrors[int64(code)]
		statusErrorGroups[statueErr.Status()] = append(statusErrorGroups[statueErr.Status()], statueErr)
	}

	p, _ := build.Import(g.pkgImportPath, "", build.ImportComment)
	buf := bytes.NewBufferString(fmt.Sprintf(`package %s

//go:generate tools gen error
import (
	"net/http"
	"github.com/johnnyeven/libtools/courier/status_error"
)
`, p.Name))

	for status, statusErrList := range statusErrorGroups {
		buf.WriteString("const (\n")

		index := 0
		for i, statusErr := range statusErrList {
			count := int(statusErr.Code) - status*1e3
			firstLine := statusErr.Msg
			if statusErr.CanBeErrorTalk {
				firstLine = "@errTalk " + firstLine
			}
			comments := fmt.Sprintf(`
// %s
// %s`, firstLine, statusErr.Desc)

			if i == 0 {
				index = count - i
				buf.WriteString(fmt.Sprintf(`%s
%s status_error.StatusErrorCode = %s * 1e3 + iota + %d
`, comments, statusErr.Key, httpCode(status), count),
				)
			} else {
				index++
				for count > index {
					buf.WriteString(fmt.Sprintf("_%d\n", index))
					index++
				}
				buf.WriteString(fmt.Sprintf(`%s
%s
`, comments, statusErr.Key))
			}
		}

		buf.WriteString(")\n\n")
	}

	if g.DryRun {
		fmt.Println(buf.String())
	} else {
		dir, _ := filepath.Rel(cwd, p.Dir)
		outputs.Add(path.Join(dir, "status_err_codes.go"), buf.String())
	}
	return outputs
}

var httpCodes = map[int]string{
	http.StatusBadRequest:          "http.StatusBadRequest",
	http.StatusNotFound:            "http.StatusNotFound",
	http.StatusForbidden:           "http.StatusForbidden",
	http.StatusTooManyRequests:     "http.StatusTooManyRequests",
	http.StatusConflict:            "http.StatusConflict",
	http.StatusInternalServerError: "http.StatusInternalServerError",
}

func httpCode(code int) string {
	if httpCode, ok := httpCodes[code]; ok {
		return httpCode
	}
	return fmt.Sprintf("%d", code)
}
