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

// getProjectsCmd represents the getProjects command
var getProjectsCmd = &cobra.Command{
	Use:   "getProjects",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires an Org Id")
		}

		client := tools.NewHttpclient(config.Instance)
		if debug {
			client.Debug = true
		}

		prjs := domain.NewProjects(client, args[0])

		var ret string
		var err error

		if rawOuput {
			ret, err = prjs.GetRaw()
			if err != nil {
				return err
			}
			fmt.Println(ret)
			return nil
		}

		err = prjs.Get()
		if err != nil {
			return err
		}

		if quiet {
			ret, _ = prjs.Quiet()
		} else if names {
			ret, _ = prjs.Names()
		} else {
			ret, _ = prjs.String()
		}
		fmt.Println(ret)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getProjectsCmd)

	getProjectsCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Prints only ids")
	getProjectsCmd.PersistentFlags().BoolVarP(&names, "names", "n", false, "Prints only names")
	getProjectsCmd.PersistentFlags().BoolVarP(&rawOuput, "raw", "r", false, "Prints raw json output from api")
}
