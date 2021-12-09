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

// getProjectsCmd represents the getProjects command
var getProjectsCmd = &cobra.Command{
	Use:   "getProjects",
	Short: "get the list of projects for an Org",
	Long: `get the list of projects for an Org. For example:
snykctl org_id 

getProject commands accepts filters such as --tag --env --lifecycle and --criticality
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)

		prjs := domain.NewProjects(client, args[0])

		var ret string
		var err error

		if checkAtLeastOneFilterSet() {
			if err := domain.ParseAttributes(attrEnvironment, attrLifecycle, attrCriticality); err != nil {
				return err
			}

			mTags, err := domain.ParseTags(attrTag)
			if err != nil {
				return err
			}
			if rawOutput {
				ret, err = prjs.GetRawFiltered(attrEnvironment, attrLifecycle, attrCriticality, mTags)
				if err != nil {
					return err
				}
				fmt.Println(ret)
				return nil
			}
			if err = prjs.GetFiltered(attrEnvironment, attrLifecycle, attrCriticality, mTags); err != nil {
				return err
			}
		} else {
			if rawOutput {
				ret, err = prjs.GetRaw()
				if err != nil {
					return err
				}
				fmt.Println(ret)
				return nil
			}

			if err = prjs.Get(); err != nil {
				return err
			}
		}

		prjs.Print(quiet, names, verbose)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getProjectsCmd)

	getProjectsCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Prints only ids")
	getProjectsCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Prints only names")
	getProjectsCmd.PersistentFlags().BoolVarP(&rawOutput, "raw", "r", false, "Prints raw json output from api")
	getProjectsCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Prints verbose (includes attributes and tags)")

	getProjectsCmd.PersistentFlags().StringVarP(&attrEnvironment, "env", "", "", "Filters by environment [frontend | backend | internal | external | mobile | saas | on-prem | hosted | distributed]")
	getProjectsCmd.PersistentFlags().StringVarP(&attrLifecycle, "lifecycle", "", "", "Filters by lifecycle [production | development | sandbox]")
	getProjectsCmd.PersistentFlags().StringVarP(&attrCriticality, "criticality", "", "", "Filters by criticality [critical | high | medium | low]")
	getProjectsCmd.PersistentFlags().StringSliceVarP(&attrTag, "tag", "", []string{}, "Filters by tag (key1=value1;key2=value2)")
}

func checkAtLeastOneFilterSet() bool {
	if attrEnvironment != "" || attrLifecycle != "" || attrCriticality != "" || len(attrTag) > 0 {
		return true
	}
	return false
}
