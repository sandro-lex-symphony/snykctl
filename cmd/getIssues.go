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

// getIssuesCmd represents the getIssues command
var getIssuesCmd = &cobra.Command{
	Use:   "getIssues",
	Short: "get the aggregated project issues",
	Long: `get the aggregated project issues. For example:
snykctl getIssues org_id prj_id
snykctl getIssues org_id
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)
		prjs := domain.NewProjects(client, args[0])

		var res domain.ProjectIssuesResult
		var prj_id, out string
		var err error
		if len(args) > 1 {
			res, err = prjs.GetIssues(args[1], issueType)
			if err != nil {
				return err
			}
			out = domain.FormatProjectIssues(res, prj_id)

		} else {
			prjs.Get()
			for _, prj := range prjs.Projects {
				res, err = prjs.GetIssues(prj.Id, issueType)
				if err != nil {
					return err
				}
				out += domain.FormatProjectIssues(res, prj.Id)
			}
		}

		fmt.Print(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getIssuesCmd)
	getIssuesCmd.PersistentFlags().StringVarP(&issueType, "type", "t", "", "Issue Type [license | vuln]")
}
