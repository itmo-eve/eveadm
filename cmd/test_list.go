/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

// testListCmd represents the list command
var testListCmd = &cobra.Command{
	Use:   "list",
	Short: "Run shell command with arguments in 'list' action on 'test' mode",
	Long: `
Run shell command with arguments in 'list' action on 'test' mode. For example:

eveadm test list ps x
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test list called")
		run(Timeout, args)
	},
}

func init() {
	testCmd.AddCommand(testListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
