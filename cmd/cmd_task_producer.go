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
	"strconv"
)

var (
	cmdHost       string
	cmdPort       int32
	cmdChannel    string
	cmdSubject    string
	cmdBrokerType string
	cmdData       string
	cmdRandomData bool
	cmdSendCount  int
)

// cmdTaskProducerCmd represents the cmdTaskProducer command
var cmdTaskProducerCmd = &cobra.Command{
	Use:   "producer",
	Short: "task producer",
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

		for i := 0; i < cmdSendCount; i++ {
			msg := fmt.Sprintf("%s %d", cmdData, i)
			t := &constants.Task{
				ID:      strconv.FormatInt(int64(i), 10),
				Channel: cmdChannel,
				Subject: cmdSubject,
				Data:    []byte(msg),
			}
			err = agent.SendTask(t)
			if err != nil {
				panic(err)
			}

			logrus.Info(msg, "sent")
		}
	},
}

func init() {
	cmdTaskProducerCmd.Flags().StringVarP(&cmdHost, "host", "", "localhost", "")
	cmdTaskProducerCmd.Flags().Int32VarP(&cmdPort, "port", "", 9092, "")
	cmdTaskProducerCmd.Flags().StringVarP(&cmdChannel, "channel", "", "", "")
	cmdTaskProducerCmd.Flags().StringVarP(&cmdSubject, "subject", "", "", "")
	cmdTaskProducerCmd.Flags().StringVarP(&cmdBrokerType, "type", "", "KAFKA", "")
	cmdTaskProducerCmd.Flags().BoolVarP(&cmdRandomData, "rand", "", false, "")
	cmdTaskProducerCmd.Flags().StringVarP(&cmdData, "data", "", "", "")
	cmdTaskProducerCmd.Flags().IntVarP(&cmdSendCount, "count", "", 1, "")

	cmdTaskCmd.AddCommand(cmdTaskProducerCmd)
}
