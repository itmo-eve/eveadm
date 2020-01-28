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

// rktUpdateCmd represents the update command
var rktUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Run shell command with arguments in 'update' action on 'rkt' mode",
	Long: `Run shell command with arguments in 'update' action on 'rkt' mode. For example:

eveadm rkt update ps x`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented now")
	},
}

func init() {
	rktCmd.AddCommand(rktUpdateCmd)
}
