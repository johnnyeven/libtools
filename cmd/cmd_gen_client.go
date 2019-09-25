package cmd

import (
	"github.com/spf13/cobra"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/courier/client/gen"
)

var (
	cmdGenClientFlagName    string
	cmdGenClientFlagFile    string
	cmdGenClientFlagSpecURL string
)

var cmdGenClient = &cobra.Command{
	Use:   "client",
	Short: "Generate client from swagger.json",
	Run: func(cmd *cobra.Command, args []string) {
		clientGenerator := gen.ClientGenerator{
			ServiceName: cmdGenClientFlagName,
			File:        cmdGenClientFlagFile,
			SpecURL:     cmdGenClientFlagSpecURL,
			BaseClient:  "github.com/johnnyeven/libtools/courier/client.Client",
		}
		codegen.Generate(&clientGenerator)
	},
}

func init() {
	cmdGenClient.Flags().
		StringVarP(&cmdGenClientFlagSpecURL, "spec-url", "", "", "client spec url")
	cmdGenClient.Flags().
		StringVarP(&cmdGenClientFlagName, "name", "", "", "service name")
	cmdGenClient.Flags().
		StringVarP(&cmdGenClientFlagFile, "file", "", "", "client spec file")

	cmdGen.AddCommand(cmdGenClient)
}
