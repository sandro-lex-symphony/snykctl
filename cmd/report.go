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

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)
		report, err := domain.NewReport(client, args[0])
		if err != nil {
			return err
		}

		if err := report.GetIssues(issueType); err != nil {
			return err
		}

		var out string
		if output == "csv" {
			out = report.ToCsv()
		} else if output == "html" {
			out = report.ToHtmlTable()
		} else {
			out = report.ToString()
		}

		fmt.Print(out)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.PersistentFlags().StringVarP(&issueType, "type", "t", "", "Issue Type [license | vuln]")
	reportCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output type [csv | html]")
}
