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

// addAttributesCmd represents the addAttributes command
var addAttributesCmd = &cobra.Command{
	Use:   "addAttributes",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, false)
		prjs := domain.NewProjects(client, args[0])
		err := prjs.AddAttributes(args[1], filterEnvironment, filterLifecycle, "")
		if err != nil {
			return err
		}

		fmt.Println("OK")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addAttributesCmd)
	addAttributesCmd.PersistentFlags().StringVarP(&filterEnvironment, "env", "", "", "Filters by environment (frontend | backend | internal | external | mobile | saas | on-prem | hosted | distributed)")
	addAttributesCmd.PersistentFlags().StringVarP(&filterLifecycle, "lifecycle", "", "", "Filters by lifecycle (production | development | sandbox)")

}
