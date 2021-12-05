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
	Short: "Gets the list of Snyk Organisations for the given token",
	Long: `Gets the list of Snyk Organisations for the given token
Example
snykctl getOrgs
snykctl getOrgs --quiet
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance)
		if debug {
			client.Debug = true
		}

		var out string
		var err error

		orgs := domain.NewOrgs(client)

		if rawOuput {
			out, err = orgs.GetRaw()
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		}

		err = orgs.Get()
		if err != nil {
			return err
		}

		if quiet {
			out, _ = orgs.Quiet()
		} else if names {
			out, _ = orgs.Names()
		} else {
			out, _ = orgs.String()
		}

		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getOrgsCmd)

	getOrgsCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Prints only ids")
	getOrgsCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Prints only names")
	getOrgsCmd.PersistentFlags().BoolVarP(&rawOuput, "raw", "r", false, "Prints raw json output from api")
}
