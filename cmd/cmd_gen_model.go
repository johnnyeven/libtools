package cmd

import (
	"github.com/spf13/cobra"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/sqlx/gen"
)

var cmdGenModelFlagDatabase string
var cmdGenModelFlagTableName string
var cmdGenModelFlagTableInterfaces bool
var cmdGenModelFlagWithComments bool
var cmdGenModelFlagFieldSoftDelete string
var cmdGenModelFlagConstSoftDeleteTrue string
var cmdGenModelFlagConstSoftDeleteFalse string

var cmdGenModel = &cobra.Command{
	Use:   "model",
	Short: "generate db model method",
	Run: func(cmd *cobra.Command, args []string) {
		if cmdGenModelFlagDatabase == "" {
			panic("database must be defined")
		}

		for _, structName := range args {
			clientGenerator := gen.SqlFuncGenerator{}
			clientGenerator.StructName = structName
			clientGenerator.Database = cmdGenModelFlagDatabase
			clientGenerator.TableName = cmdGenModelFlagTableName
			clientGenerator.TableName = cmdGenModelFlagTableName
			clientGenerator.WithTableInterfaces = cmdGenModelFlagTableInterfaces
			clientGenerator.WithComments = cmdGenModelFlagWithComments
			clientGenerator.FieldSoftDelete = cmdGenModelFlagFieldSoftDelete
			clientGenerator.ConstSoftDeleteTrue = cmdGenModelFlagConstSoftDeleteTrue
			clientGenerator.ConstSoftDeleteFalse = cmdGenModelFlagConstSoftDeleteFalse
			codegen.Generate(&clientGenerator)
		}
	},
}

func init() {
	cmdGenModel.Flags().
		StringVarP(&cmdGenModelFlagDatabase, "database", "", "", "(required) register model to database var")
	cmdGenModel.Flags().
		StringVarP(&cmdGenModelFlagTableName, "table-name", "t", "", "custom table name")
	cmdGenModel.Flags().
		BoolVarP(&cmdGenModelFlagTableInterfaces, "with-table-interfaces", "", true, "with table interface TableName T D")
	cmdGenModel.Flags().
		BoolVarP(&cmdGenModelFlagWithComments, "with-comments", "", false, "use comments")
	cmdGenModel.Flags().
		StringVarP(&cmdGenModelFlagFieldSoftDelete, "field-soft-delete", "", "", "custom soft delete field")
	cmdGenModel.Flags().
		StringVarP(&cmdGenModelFlagConstSoftDeleteTrue, "const-soft-delete-true", "", "", "custom soft delete value true")
	cmdGenModel.Flags().
		StringVarP(&cmdGenModelFlagConstSoftDeleteFalse, "const-soft-delete-false", "", "", "custom soft delete value false")

	cmdGen.AddCommand(cmdGenModel)
}
