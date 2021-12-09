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

// copyUsersCmd represents the copyUsers command
var copyUsersCmd = &cobra.Command{
	Use:   "copyUsers",
	Short: "copy users from one Org to another Org",
	Long: `Copy users from one Org to another Org. For example:
snykctl copyUsers org1 org2
`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := tools.NewHttpclient(config.Instance, debug)

		if err := domain.CopyUsers(client, args[0], args[1]); err != nil {
			return err
		}
		fmt.Println("OK")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyUsersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyUsersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyUsersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
