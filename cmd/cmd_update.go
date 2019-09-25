package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/johnnyeven/libtools/executil"
	"github.com/johnnyeven/libtools/project/govendor"
)

var cmdUpdateFlagUpdate bool

var cmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "go vendor update",
	Run: func(cmd *cobra.Command, args []string) {
		vendor, err := govendor.LoadGoVendorJSON("./vendor/vendor.json")
		if err == nil {
			importPathList := vendor.ListImportPath()
			govendor.UpdatePkgs(importPathList...)
			if cmdUpdateFlagUpdate {
				executil.StdRun(exec.Command("govendor", "update", "+v"))
			}
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdUpdate)

	cmdUpdate.Flags().
		BoolVarP(&cmdUpdateFlagUpdate, "update", "u", false, "update vendor")
}
