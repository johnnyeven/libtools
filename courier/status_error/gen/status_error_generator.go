package gen

import (
	"bytes"
	"fmt"
	"go/build"
	"go/parser"
	"go/types"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/codegen/loaderx"
	"github.com/johnnyeven/libtools/courier/status_error"
)

type StatusErrorGenerator struct {
	pkgImportPath    string
	program          *loader.Program
	statusErrorCodes map[*types.Package]status_error.StatusErrorCodeMap
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
	g.statusErrorCodes = map[*types.Package]status_error.StatusErrorCodeMap{}
}

func (g *StatusErrorGenerator) Pick() {
	statusErrorCodeType := reflect.TypeOf(status_error.StatusErrorCode(0))
	statusErrorCodeTypeFullName := fmt.Sprintf("%s.%s", statusErrorCodeType.PkgPath(), statusErrorCodeType.Name())

	for pkg, pkgInfo := range g.program.AllPackages {
		if pkg.Path() != g.pkgImportPath {
			continue
		}
		for ident, obj := range pkgInfo.Defs {
			if constObj, ok := obj.(*types.Const); ok {
				if strings.HasSuffix(constObj.Type().String(), statusErrorCodeTypeFullName) {
					key := constObj.Name()
					if key == "_" {
						continue
					}

					doc := loaderx.CommentsOf(g.program.Fset, ident, pkgInfo.Files...)
					code, _ := strconv.ParseInt(constObj.Val().String(), 10, 64)
					msg, desc, canBeErrTalk := ParseStatusErrorDesc(doc)

					if g.statusErrorCodes[pkg] == nil {
						g.statusErrorCodes[pkg] = status_error.StatusErrorCodeMap{}
					}
					g.statusErrorCodes[pkg].Register(key, code, msg, desc, canBeErrTalk)
				}
			}
		}
	}
}

func (g *StatusErrorGenerator) Output(cwd string) codegen.Outputs {
	statusErrorCodeType := reflect.TypeOf(status_error.StatusErrorCode(0))
	outputs := codegen.Outputs{}
	for pkg, statusErrorCodeMap := range g.statusErrorCodes {
		p, _ := build.Import(pkg.Path(), "", build.FindOnly)
		dir, _ := filepath.Rel(cwd, p.Dir)
		content := fmt.Sprintf(`
			package %s

			import(
				%s
			)

			%s `,
			pkg.Name(),
			strconv.Quote(statusErrorCodeType.PkgPath()),
			g.toRegisterInit(statusErrorCodeMap),
		)
		outputs.Add(codegen.GeneratedSuffix(path.Join(dir, "status_err_codes.go")), content)
	}
	return outputs
}

func (g *StatusErrorGenerator) toRegisterInit(statusErrorCodeMap status_error.StatusErrorCodeMap) string {
	buffer := bytes.Buffer{}
	buffer.WriteString("func init () {")

	registerMethod := "status_error.StatusErrorCodes"
	pkgs := strings.Split(g.pkgImportPath, "/")
	if strings.HasPrefix(registerMethod, pkgs[len(pkgs)-1]) {
		registerMethod = "StatusErrorCodes"
	}

	statusErrorCodeList := []int{}
	for _, statusErrorCode := range statusErrorCodeMap {
		statusErrorCodeList = append(statusErrorCodeList, int(statusErrorCode.Code))
	}
	sort.Ints(statusErrorCodeList)

	for _, code := range statusErrorCodeList {
		statusErrorCode := statusErrorCodeMap[int64(code)]

		buffer.WriteString(fmt.Sprintf(
			"%s.Register(%s, %d, %s, %s, %s)\n",
			registerMethod,
			strconv.Quote(statusErrorCode.Key),
			statusErrorCode.Code,
			strconv.Quote(statusErrorCode.Msg),
			strconv.Quote(statusErrorCode.Desc),
			strconv.FormatBool(statusErrorCode.CanBeErrorTalk),
		))
	}

	buffer.WriteString("}")
	return buffer.String()
}

func ParseStatusErrorDesc(str string) (msg string, desc string, canBeErrTalk bool) {
	lines := strings.Split(str, "\n")
	firstLine := strings.Split(lines[0], "@errTalk")

	if len(firstLine) > 1 {
		canBeErrTalk = true
		msg = strings.TrimSpace(firstLine[1])
	} else {
		canBeErrTalk = false
		msg = strings.TrimSpace(firstLine[0])
	}

	if len(lines) > 1 {
		desc = strings.TrimSpace(strings.Join(lines[1:], "\n"))
	}
	return
}
