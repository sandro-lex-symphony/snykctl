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
	"bufio"
	"fmt"
	"os"

	"snykctl/internal/config"

	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configure()
	},
}

func configure() {
	var config config.ConfigProperties
	config.Sync()

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("token [.... %s]: ", config.ObfuscatedToken())
	text, _ := reader.ReadString('\n')
	if len(text) > 2 {
		config.SetToken(text[:len(text)-1])
	}

	fmt.Printf("group_id [.... %s]: ", config.ObfuscatedId())
	text, _ = reader.ReadString('\n')
	if len(text) > 2 {
		config.SetId(text[:len(text)-1])
	}

	fmt.Printf("timeout [%d]: ", config.Timeout())
	text, _ = reader.ReadString('\n')
	if len(text) > 1 {
		config.SetTimeoutStr(text[:len(text)-1])
	}

	fmt.Printf("worker size [%d]: ", config.WorkerSize())
	text, _ = reader.ReadString('\n')
	if len(text) > 1 {
		config.SetWorkerSizeStr(text[:len(text)-1])
	}

	config.WriteConf()
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
