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

// getGroupUsersCmd represents the getGroupUsers command
var getGroupUsersCmd = &cobra.Command{
	Use:   "getGroupUsers",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var group_id string
		if len(args) != 1 {
			// get the group id from conf
			group_id = config.Instance.Id()
		} else {
			group_id = args[0]
		}

		client := tools.NewHttpclient(config.Instance, debug)

		users := domain.NewUsers(client, group_id)

		var out string
		var err error

		if rawOutput {
			out, err = users.GetGroupRaw()
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		}

		err = users.GetGroup()
		if err != nil {
			return err
		}

		if quiet {
			out = users.Quiet()
		} else if names {
			out = users.Name()
		} else {
			out = users.String()
		}
		fmt.Print(out)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getGroupUsersCmd)

	getGroupUsersCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Prints only ids")
	getGroupUsersCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Prints only names")
	getGroupUsersCmd.PersistentFlags().BoolVarP(&rawOutput, "raw", "r", false, "Prints raw json output from api")
}