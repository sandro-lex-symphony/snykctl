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

// getProjectCmd represents the getProject command
var getProjectCmd = &cobra.Command{
	Use:   "getProject",
	Short: "print the raw json with the project values",
	Long: `print the raw json with the project values For example:
snykctl org_id prj_id
`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)

		prjs := domain.NewProjects(client, args[0])

		var ret string
		var err error

		ret, err = prjs.GetRawProject(args[1])
		if err != nil {
			return err
		}
		fmt.Println(ret)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getProjectCmd)
}
