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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteAllProjectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteAllProjectsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
