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

// getProjectIssuesCmd represents the getProjectIssues command
var getProjectIssuesCmd = &cobra.Command{
	Use:   "getProjectIssues",
	Short: "get the aggregated project issues",
	Long: `get the aggregated project issues. For example:
snykctl getProjectIssues org_id prj_id
`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)
		prjs := domain.NewProjects(client, args[0])

		out, err := prjs.GetProjectIssues(args[1])
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getProjectIssuesCmd)
}
