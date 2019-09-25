package cmd

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/johnnyeven/libtools/project/commit_msg"
)

var cmdCommitMsgFlagInit bool

var cmdCommitMsg = &cobra.Command{
	Use:   "commit_msg",
	Short: "commit_msg hooks",
	Run: func(cmd *cobra.Command, args []string) {
		if cmdCommitMsgFlagInit {
			ioutil.WriteFile(".git/hooks/commit-msg", []byte(`#!/bin/sh
tools commit_msg $1
`), os.ModePerm)
		} else {
			data, _ := ioutil.ReadFile(args[0])
			err := commit_msg.CheckCommit(string(data))
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	cmdCommitMsg.Flags().
		BoolVarP(&cmdCommitMsgFlagInit, "init", "", false, "init commit msg hook")

	cmdRoot.AddCommand(cmdCommitMsg)
}
