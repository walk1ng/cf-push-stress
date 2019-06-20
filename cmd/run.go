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
	"github.com/walk1ng/cf-push-stress/model"
	"github.com/walk1ng/cf-push-stress/utils"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the stress",
	Run: func(cmd *cobra.Command, args []string) {

		var testResults []model.PushAppResult

		pushCount, _ := cmd.Flags().GetInt("requests")
		cfHost, _ := cmd.Flags().GetString("host")

		if mode, _ := cmd.Flags().GetBool("serial"); mode {
			testResults = utils.SerialPush(pushCount, cfHost)
		}

		if mode, _ := cmd.Flags().GetBool("concurrency"); mode {
			conc, _ := cmd.Flags().GetInt("conc")
			rounds := pushCount / conc
			for rd := 1; rd <= rounds; rd++ {
				tr := utils.ConcurrencyPush(rd, conc, cfHost)
				testResults = append(testResults, tr...)
			}
		}

		utils.SummaryTest(pushCount, testResults)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.PersistentFlags().String("host", "", "cf host domain")
	runCmd.PersistentFlags().Bool("serial", false, "serially push")
	runCmd.PersistentFlags().Bool("concurrency", false, "concurrency push")
	runCmd.PersistentFlags().Int("requests", 5, "number of push to perform")
	runCmd.PersistentFlags().Int("conc", 5, "number of multiple push to make at a time ")
}
