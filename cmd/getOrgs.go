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

// getOrgsCmd represents the getOrgs command
var getOrgsCmd = &cobra.Command{
	Use:   "getOrgs",
	Short: "gets the list of Snyk Organisations for the given token",
	Long: `gets the list of Snyk Organisations for the given token
Example
snykctl getOrgs
snykctl getOrgs --quiet
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)

		var out string
		var err error

		orgs := domain.NewOrgs(client)

		if rawOutput {
			out, err = orgs.GetRaw()
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		}

		if err = orgs.Get(); err != nil {
			return err
		}
		orgs.Print(quiet, names)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getOrgsCmd)

	getOrgsCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Prints only ids")
	getOrgsCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Prints only names")
	getOrgsCmd.PersistentFlags().BoolVarP(&rawOutput, "raw", "r", false, "Prints raw json output from api")
}
