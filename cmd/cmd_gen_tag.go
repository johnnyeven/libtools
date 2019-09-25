package cmd

import (
	"github.com/spf13/cobra"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/sqlx/gen"
)

var withDefaults bool

var cmdGenTag = &cobra.Command{
	Use:   "tag",
	Short: "generate db model tags",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}
		g := gen.TagGenerator{
			WithDefaults: withDefaults,
		}
		g.StructNames = args
		codegen.Generate(&g)
	},
}

func init() {
	cmdGenTag.Flags().
		BoolVarP(&withDefaults, "defaults", "", false, "generate tags with Default")

	cmdGen.AddCommand(cmdGenTag)
}
