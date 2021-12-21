/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"snykctl/internal/config"
	"snykctl/internal/domain"
	"snykctl/internal/tools"

	"github.com/spf13/cobra"
)

// countIssuesCmd represents the countIssues command
var countIssuesCmd = &cobra.Command{
	Use:   "countIssues",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)

		result, err := domain.IssueCount(client, args[0])
		if err != nil {
			return err
		}

		out := domain.FormatIssueCountResults(result, output)
		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(countIssuesCmd)
	countIssuesCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output type [html]")
}
