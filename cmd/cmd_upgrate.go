package cmd

import (
	"go/build"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade tools",
	Run: func(cmd *cobra.Command, args []string) {
		pkg, _ := build.Import("golib/tools", "", build.ImportComment)
		os.Chdir(pkg.ImportComment)
		cwd, _ := os.Getwd()

		logrus.Infof("upgrading golib/tools in %s", cwd)

		{
			output, err := exec.Command("git", "pull", "--rebase").CombinedOutput()
			if err != nil {
				logrus.Panicf(err.Error())
			}
			logrus.Info(string(output))
		}

		{
			_, err := exec.Command("go", "install").CombinedOutput()
			if err != nil {
				logrus.Panicf(err.Error())
			}
			logrus.Infof("upgrade success")
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdUpgrade)
}
