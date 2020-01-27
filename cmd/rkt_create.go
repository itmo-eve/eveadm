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
	Use:   "create",
	Short: "Run shell command with arguments in 'create' action on 'rkt' mode",
	Long: `
Run shell command with arguments in 'create' action on 'rkt' mode. For example:

eveadm rkt create
`,
	Run: func(cmd *cobra.Command, args []string) {
		err, args, envs := rktCreateToCmd(rktctx)
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
	rktCmd.AddCommand(rktCreateCmd)
	rktCreateCmd.Flags().StringVar(&rktctx.uuidFile, "uuid-file-save", "", "File to save uuid")
	rktCreateCmd.Flags().StringVar(&rktctx.imageUUID, "image-hash", "", "Hash of image")
	rktCreateCmd.Flags().StringVar(&rktctx.xenCfgFilename, "xen-cfg-filename", "", "File with xen cfg for stage1")
	rktCreateCmd.Flags().StringVar(&rktctx.stage1Path, "stage1-path", "/usr/sbin/stage1-xen.aci", "Stage1 path")

	//Workaround to start in Ubuntu
	rktCreateCmd.Flags().BoolVar(&rktctx.noOverlay, "no-overlay", false, "Run without overlay")

	rktCreateCmd.Flags().BoolVar(&rktctx.runPaused, "paused", true, "Run paused")
	err := rktCreateCmd.MarkFlagRequired("image-hash")
	if err != nil {
		log.Fatalf("Failed to mark required flag")
	}
}
