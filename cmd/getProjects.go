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
	Short: "Get the list of projects in the Org",
	Long: `Prints the list of projetcs in the org

Example:
  snykctl getProjects org_id [flags]
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)

		prjs := domain.NewProjects(client, args[0])

		var ret string
		var err error

		if rawOutput {
			ret, err = prjs.GetRaw()
			if err != nil {
				return err
			}
			fmt.Println(ret)
			return nil
		}

		if checkAtLeastOneFilterSet() {
			err = parseFilters()
			if err != nil {
				return err
			}
			err = prjs.GetFiltered(filterEnvironment, filterLifecycle)
			if err != nil {
				return err
			}
		} else {
			err = prjs.Get()
			if err != nil {
				return err
			}
		}

		if quiet {
			ret, _ = prjs.Quiet()
		} else if names {
			ret, _ = prjs.Names()
		} else {
			ret, _ = prjs.String()
		}
		fmt.Println(ret)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getProjectsCmd)

	getProjectsCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Prints only ids")
	getProjectsCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Prints only names")
	getProjectsCmd.PersistentFlags().BoolVarP(&rawOutput, "raw", "r", false, "Prints raw json output from api")

	getProjectsCmd.PersistentFlags().StringVarP(&filterEnvironment, "env", "", "", "Filters by environment (frontend | backend | internal | external | mobile | saas | on-prem | hosted | distributed)")
	getProjectsCmd.PersistentFlags().StringVarP(&filterLifecycle, "lifecycle", "", "", "Filters by lifecycle (production | development | sandbox)")
}

func checkAtLeastOneFilterSet() bool {
	// if (Key != "" && Value != "") || FilterEnvironment != "" || FilterLifecycle != "" {
	if filterEnvironment != "" || filterLifecycle != "" {
		return true
	}
	return false
}

func parseFilters() error {
	if filterEnvironment != "" {
		if !tools.Contains(validEnvironments[:], filterEnvironment) {
			return fmt.Errorf("invalid environment value: %s\nValid values: %v", filterEnvironment, validEnvironments[:])
		}
	}

	if filterLifecycle != "" {
		if !tools.Contains(validLifecycle[:], filterLifecycle) {
			return fmt.Errorf("invalid lifecycle value: %s\nValid values: %v", filterLifecycle, validLifecycle[:])
		}
	}

	return nil
}
