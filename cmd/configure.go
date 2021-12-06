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
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("token [.... %s]: ", config.Instance.ObfuscatedToken())
	text, _ := reader.ReadString('\n')
	if len(text) > 2 {
		config.Instance.SetToken(text[:len(text)-1])
	}

	fmt.Printf("group_id [.... %s]: ", config.Instance.ObfuscatedId())
	text, _ = reader.ReadString('\n')
	if len(text) > 2 {
		config.Instance.SetId(text[:len(text)-1])
	}

	fmt.Printf("timeout [%d]: ", config.Instance.Timeout())
	text, _ = reader.ReadString('\n')
	if len(text) > 1 {
		config.Instance.SetTimeoutStr(text[:len(text)-1])
	}

	fmt.Printf("worker size [%d]: ", config.Instance.WorkerSize())
	text, _ = reader.ReadString('\n')
	if len(text) > 1 {
		config.Instance.SetWorkerSizeStr(text[:len(text)-1])
	}

	config.Instance.WriteConf()
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
