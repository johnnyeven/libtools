package govendor

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/johnnyeven/libtools/executil"
)

var pkgsNoGoGet = []string{
	"g7pay/",
	"devops/",
	"golib/",
}

func isLocalPkg(pkgName string) bool {
	for _, pkgPrefix := range pkgsNoGoGet {
		if strings.HasPrefix(pkgName, pkgPrefix) {
			return true
		}
	}
	return false
}

func UpdatePkgs(importPaths ...string) {
	sort.Strings(importPaths)
	needToUpdates := map[string]bool{}

	for _, importPath := range importPaths {
		isSubPkg := false
		for p := range needToUpdates {
			if strings.HasPrefix(importPath, p) {
				isSubPkg = true
			}
		}
		if !isSubPkg {
			needToUpdates[importPath] = true
		}
	}

	for importPath := range needToUpdates {
		UpdatePkg(importPath)
	}
}

func UpdatePkg(importPath string) {
	pkg, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		if !isLocalPkg(importPath) {
			executil.StdRun(exec.Command("go", "get", "-v", importPath))
			return
		}
		goPath := os.Getenv("GOPATH")
		os.Chdir(path.Join(strings.Split(goPath, ":")[0], "src"))

		gitRepo := fmt.Sprintf(`git@git.chinawayltd.com:%s.git`, importPath)
		executil.StdRun(exec.Command("git", "clone", gitRepo, importPath))
	} else {
		os.Chdir(pkg.Dir)
		executil.StdRun(exec.Command("git", "pull", "--rebase"))
	}
}
