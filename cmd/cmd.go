package cmd

import (
	"fmt"
	"go/build"
	"os"

	"github.com/spf13/cobra"
)

var (
	packageName string
)

func getPackageName() string {
	pwd, _ := os.Getwd()
	pkg, err := build.ImportDir(pwd, build.FindOnly)
	if err != nil {
		panic(err)
	}
	return pkg.ImportPath
}

var cmdRoot = &cobra.Command{
	Use:   "libtools",
	Short: "profzone tools",
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cmdRoot.PersistentFlags().StringVarP(&packageName, "package", "p", getPackageName(), "package name for generating")
}
