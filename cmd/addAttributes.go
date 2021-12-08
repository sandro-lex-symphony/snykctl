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
	Short: "add attributes to projects",
	Long: `attributes are static and non-configurable fields which allow to add additional metadata to a project. 
	Attributes have a pre-defined list of values that a user can select from.
	--env [frontend | backend | internal ...]
	--lifecycle [production | development | sandbox ]
	--criticality [critical | high | medium | low ]

	snykctl addAttributes --env frontend org_id prj_id
`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, false)
		prjs := domain.NewProjects(client, args[0])
		err := prjs.AddAttributes(args[1], attrEnvironment, attrLifecycle, attrCriticality)
		if err != nil {
			return err
		}

		fmt.Println("OK")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addAttributesCmd)
	addAttributesCmd.PersistentFlags().StringVarP(&attrEnvironment, "env", "", "", "Environment (frontend | backend | internal | external | mobile | saas | on-prem | hosted | distributed)")
	addAttributesCmd.PersistentFlags().StringVarP(&attrLifecycle, "lifecycle", "", "", "Lifecycle (production | development | sandbox)")
	addAttributesCmd.PersistentFlags().StringVarP(&attrCriticality, "criticality", "", "", "Criticality [critical | high | medium | low]")

}
