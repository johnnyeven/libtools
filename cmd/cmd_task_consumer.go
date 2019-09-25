// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/johnnyeven/libtools/task"
	"github.com/johnnyeven/libtools/task/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

// cmdTaskConsumerCmd represents the cmdTaskConsumer command
var cmdTaskConsumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "task consumer",
	Run: func(cmd *cobra.Command, args []string) {
		brokerType, err := constants.ParseBrokerTypeFromString(cmdBrokerType)
		if err != nil {
			panic(err)
		}
		agent := &task.Agent{
			ConnectionInfo: constants.ConnectionInfo{
				Host: cmdHost,
				Port: cmdPort,
			},
			Channel:    cmdChannel,
			BrokerType: brokerType,
		}

		agent.Init()
		agent.Register(cmdSubject, func(i2 *constants.Task) (i interface{}, e error) {
			fmt.Println(string(i2.Data))
			return nil, nil
		})
		agent.Start(1)
		logrus.Info("worker started")

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)

		<-sig
		logrus.Info("worker stopped")
	},
}

func init() {
	cmdTaskConsumerCmd.Flags().StringVarP(&cmdHost, "host", "", "localhost", "")
	cmdTaskConsumerCmd.Flags().Int32VarP(&cmdPort, "port", "", 9092, "")
	cmdTaskConsumerCmd.Flags().StringVarP(&cmdChannel, "channel", "", "", "")
	cmdTaskConsumerCmd.Flags().StringVarP(&cmdSubject, "subject", "", "", "")
	cmdTaskConsumerCmd.Flags().StringVarP(&cmdBrokerType, "type", "", "KAFKA", "")

	cmdTaskCmd.AddCommand(cmdTaskConsumerCmd)
}
