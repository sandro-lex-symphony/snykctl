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
	"errors"
	"fmt"
	"snykctl/internal/config"
	"snykctl/internal/domain"
	"snykctl/internal/tools"

	"github.com/spf13/cobra"
)

// getOrgNameCmd represents the getOrgName command
var getOrgNameCmd = &cobra.Command{
	Use:   "getOrgName",
	Short: "Returns the name of the gievn org",
	Long: `Returns the name of the given org.  For example:
snykctl getOrgName [org-id]`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires an Org Id")
		}
		client := tools.NewHttpclient(config.Instance)
		if debug {
			client.Debug = true
		}

		orgs := domain.NewOrgs(client)
		out, err := orgs.GetOrgName(args[0])
		if err != nil {
			return err
		}
		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getOrgNameCmd)
}
