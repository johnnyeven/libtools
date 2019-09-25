package cmd

import (
	"github.com/spf13/cobra"

	"github.com/johnnyeven/libtools/gen_method"
)

var cmdGenMethodFlagNoTableName bool

var cmdGenMethod = &cobra.Command{
	Use:   "method",
	Short: "generate db model method",
	Run: func(cmd *cobra.Command, args []string) {
		eg := gen_method.NewScanner(packageName)
		if args[0] != "" {
			eg.Output(args[0], cmdGenMethodFlagNoTableName)
		}
	},
}

func init() {
	cmdGenMethod.Flags().
		BoolVarP(&cmdGenMethodFlagNoTableName, "no-table-name", "", false, "skip tableName")

	cmdGen.AddCommand(cmdGenMethod)
}
