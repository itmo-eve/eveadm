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
	"log"
)

// rktStartCmd represents the start command
var rktStartCmd = &cobra.Command{
	Use:   "start id",
	Short: "Run shell command with arguments in 'start' action on 'rkt' mode",
	Long: `Run shell command with arguments in 'start' action on 'rkt' mode. For example:

eveadm rkt start uuid
`, Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		if rktctx.stage1Type == "xen" {
			rktctx.containerUUID = arg
			err, args, envs := rktctx.rktStartToCmd()
			if err != nil {
				log.Fatalf("Error in obtain params in %s", cmd.Name())
			}
			Run(cmd, Timeout, args, envs)
		} else {
			fmt.Println("Not implemented for common type of stage1")
		}
	},
}

func init() {
	rktCmd.AddCommand(rktStartCmd)
}
