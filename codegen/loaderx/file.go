package loaderx

import (
	"go/ast"
	"go/build"
	"os"
)

func FileOf(targetNode ast.Node, files ...*ast.File) *ast.File {
	for _, file := range files {
		if file.Pos() <= targetNode.Pos() && file.End() > targetNode.Pos() {
			return file
		}
	}
	return nil
}

func ResolvePkgImport(pkgImportPath string) string {
	cwd, _ := os.Getwd()
	pkg, _ := build.Import(pkgImportPath, cwd, build.FindOnly)
	return pkg.ImportPath
}
