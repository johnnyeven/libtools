package loaderx

import (
	"fmt"
	"go/parser"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/tools/go/loader"

	"github.com/johnnyeven/libtools/codegen"
)

func replaceTagPlaceholder(s string) string {
	return strings.Replace(s, "^^", "`", -1)
}

func NewTestProgram(content string) (pkgImportPath string, program *loader.Program) {
	ldr := loader.Config{
		ParserMode: parser.ParseComments,
	}

	cwd, _ := os.Getwd()
	filename := fmt.Sprintf("%s/%s/main.go", cwd, uuid.New().String())
	defer codegen.CreateTempFile(filename, replaceTagPlaceholder(content))()

	pkgImportPath = codegen.GetPackageImportPath(filepath.Dir(filename))
	ldr.Import(pkgImportPath)

	p, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	return pkgImportPath, p
}
