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

// getIgnoresCmd represents the getIgnores command
var getIgnoresCmd = &cobra.Command{
	Use:   "getIgnores",
	Short: "get the list of ignores for a project",
	Long: `get the list of ignores for a project. For example:
snykctl getIgnores org_id prj_id
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)
		var out string
		if len(args) > 1 {
			ignores, err := domain.GetProjectIgnores(client, args[0], args[1])
			if err != nil {
				return err
			}
			out = domain.FormatIgnoreResult(ignores, "")
		} else {
			// get all org ignores
			prjs := domain.NewProjects(client, args[0])
			if err := prjs.Get(); err != nil {
				return err
			}
			for _, prj := range prjs.Projects {
				ignores, err := domain.GetProjectIgnores(client, args[0], prj.Id)
				if err != nil {
					return err
				}
				out += domain.FormatIgnoreResult(ignores, prj.Id)
			}

		}
		fmt.Print(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getIgnoresCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getIgnoresCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getIgnoresCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
