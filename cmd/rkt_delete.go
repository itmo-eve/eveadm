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

// rktDeleteCmd represents the delete command
var rktDeleteCmd = &cobra.Command{
	Use:   "delete uuid",
	Short: "Run shell command with arguments in 'delete' action on 'rkt' mode",
	Long: `
Run shell command with arguments in 'delete' action on 'rkt' mode. For example:

eveadm rkt delete ps x`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		isImage, err := cmd.Flags().GetBool("image")
		if err != nil {
			log.Fatalf("Error in get param image in %s", cmd.Name())
		}
		var envs string
		if isImage {
			rktctx.imageUUID = arg
			err, args, envs = rktDeleteImageToCmd(rktctx)
		} else {
			rktctx.containerUUID = arg
			err, args, envs = rktDeleteToCmd(rktctx)
		}
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
				_, err = fmt.Fprintf(os.Stderr, "%s", stderr.String())
				if err != nil {
					fmt.Printf("%s", stderr.String())
				}
				os.Exit(waitStatus.ExitStatus())
			} else {
				_, err = fmt.Fprintf(os.Stderr, "Execute error in %s: %s\n", cmd.Name(), err.Error())
				if err != nil {
					fmt.Printf("Execute error in %s: %s\n", cmd.Name(), err.Error())
				}
			}
		}
		fmt.Printf("%s", stdout.String())
	},
}

func init() {
	rktCmd.AddCommand(rktDeleteCmd)
	rktDeleteCmd.Flags().BoolP("image", "i", false, "Work with images")
}
