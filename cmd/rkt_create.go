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

// rktCreateCmd represents the create command
var rktCreateCmd = &cobra.Command{
	Use:   "create url/uuid",
	Short: "Run shell command with arguments in 'create' action on 'rkt' mode",
	Long: `
Run shell command with arguments in 'create' action on 'rkt' mode. For example:

eveadm rkt create
`, Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		isImage, err := cmd.Flags().GetBool("image")
		if err != nil {
			log.Fatalf("Error in get param image in %s", cmd.Name())
		}
		var envs string
		if isImage {
			rktctx.imageUrl = arg
			err, args, envs = rktCreateImageToCmd(rktctx)
		} else {
			rktctx.imageUUID = arg
			err, args, envs = rktCreateToCmd(rktctx)
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
	rktCmd.AddCommand(rktCreateCmd)
	rktCreateCmd.Flags().BoolP("image", "i", false, "Work with images")
	rktCreateCmd.Flags().StringVar(&rktctx.uuidFile, "uuid-file-save", "", "File to save uuid")
	rktCreateCmd.Flags().StringVar(&rktctx.xenCfgFilename, "xen-cfg-filename", "", "File with xen cfg for stage1")
	rktCreateCmd.Flags().StringVar(&rktctx.stage1Path, "stage1-path", "/usr/sbin/stage1-xen.aci", "Stage1 path")

	//Workaround to start in Ubuntu
	rktCreateCmd.Flags().BoolVar(&rktctx.noOverlay, "no-overlay", false, "Run without overlay")

	rktCreateCmd.Flags().BoolVar(&rktctx.runPaused, "paused", true, "Run paused")
}
