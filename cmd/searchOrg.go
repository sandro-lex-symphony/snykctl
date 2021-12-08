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
	"snykctl/internal/config"
	"snykctl/internal/domain"
	"snykctl/internal/tools"
	"strings"

	"github.com/spf13/cobra"
)

// searchOrgCmd represents the searchOrg command
var searchOrgCmd = &cobra.Command{
	Use:   "searchOrg",
	Short: "search Org using name",
	Long: `search Org using name For example:
snykctl searchOrg [search-term]
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)

		orgs := domain.NewOrgs(client)
		err := orgs.Get()
		if err != nil {
			return err
		}

		var filteredOrgs domain.Orgs
		for _, org := range orgs.Orgs {
			if strings.Contains(strings.ToLower(org.Name), strings.ToLower(args[0])) {
				filteredOrgs.Orgs = append(filteredOrgs.Orgs, org)
			}
		}

		filteredOrgs.Print(quiet, names)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchOrgCmd)

	searchOrgCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Prints only ids")
	searchOrgCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Prints only names")
}
