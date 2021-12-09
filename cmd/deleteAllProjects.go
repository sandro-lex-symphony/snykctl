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

// deleteAllProjectsCmd represents the deleteAllProjects command
var deleteAllProjectsCmd = &cobra.Command{
	Use:   "deleteAllProjects",
	Short: "delete all projects in an Org",
	Long: `delete all projects in an Org. For example:
snykctl deleteAllProjects org_id
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, false)
		prjs := domain.NewProjects(client, args[0])

		if checkAtLeastOneFilterSet() {
			if err := domain.ParseAttributes(attrEnvironment, attrLifecycle, attrCriticality); err != nil {
				return err
			}

			mTags, err := domain.ParseTags(attrTag)
			if err != nil {
				return err
			}

			if err = prjs.GetFiltered(attrEnvironment, attrLifecycle, attrCriticality, mTags); err != nil {
				return err
			}
		} else {
			if err := prjs.Get(); err != nil {
				return err
			}
		}

		out, err := prjs.DeleteAllProjects()
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteAllProjectsCmd)

	deleteAllProjectsCmd.PersistentFlags().StringVarP(&attrEnvironment, "env", "", "", "Filters by environment [frontend | backend | internal | external | mobile | saas | on-prem | hosted | distributed]")
	deleteAllProjectsCmd.PersistentFlags().StringVarP(&attrLifecycle, "lifecycle", "", "", "Filters by lifecycle [production | development | sandbox]")
	deleteAllProjectsCmd.PersistentFlags().StringVarP(&attrCriticality, "criticality", "", "", "Filters by criticality [critical | high | medium | low]")
	deleteAllProjectsCmd.PersistentFlags().StringSliceVarP(&attrTag, "tag", "", []string{}, "Filters by tag (key1=value1;key2=value2)")

}
