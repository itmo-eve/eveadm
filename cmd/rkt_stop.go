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
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

// rktStopCmd represents the stop command
var rktStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Run shell command with arguments in 'stop' action on 'rkt' mode",
	Long: `Run shell command with arguments in 'stop' action on 'rkt' mode. For example:

eveadm rkt stop`,
	Run: func(cmd *cobra.Command, args []string) {
		err, args, envs := rktStopToCmd(rktctx)
		if err != nil {
			log.Fatalf("Error in obtain params in %s", cmd.Name())
		}
		err, cerr, stdout, stderr := rune(Timeout, args, envs)
		if cerr != nil {
			log.Fatalf("Context error in %s", cmd.Name())
		}
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus := exitError.Sys().(syscall.WaitStatus)
				fmt.Printf("%s", stdout.String())
				fmt.Printf("%s", stderr.String())
				os.Exit(waitStatus.ExitStatus())
			} else {
				log.Fatalf("Execute error in %s", cmd.Name())
			}
		}
		fmt.Printf("%s", stdout.String())
	},
}

func init() {
	rktCmd.AddCommand(rktStopCmd)
	rktStopCmd.Flags().StringVar(&rktctx.containerUUID, "container-uuid", "", "UUID of container")
	err := rktStopCmd.MarkFlagRequired("container-uuid")
	if err != nil {
		log.Fatalf("Failed to mark required flag")
	}
}
