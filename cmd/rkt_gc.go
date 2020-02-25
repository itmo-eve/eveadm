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
	"github.com/spf13/cobra"
	"log"
)

// rktDeleteCmd represents the delete command
var rktGCCmd = &cobra.Command{
	Use:   "gc",
	Short: "Run shell command with arguments in 'gc' action on 'rkt delete' mode",
	Long: `
Run shell command with arguments in 'gc' action on 'rkt delete' mode. For example:

eveadm rkt delete gc`,
	Run: func(cmd *cobra.Command, args []string) {
		isImage, err := cmd.Flags().GetBool("image")
		if err != nil {
			log.Fatalf("Error in get param image in %s", cmd.Name())
		}
		var envs string
		if isImage {
			err, args, envs = rktctx.rktDeleteGCImageToCmd()
		} else {
			err, args, envs = rktctx.rktDeleteGCToCmd()
		}
		if err != nil {
			log.Fatalf("Error in obtain params in %s", cmd.Name())
		}
		Run(cmd, Timeout, args, envs)
	},
}

func init() {
	rktDeleteCmd.AddCommand(rktGCCmd)
	rktGCCmd.Flags().BoolP("image", "i", false, "Work with images")
	rktGCCmd.Flags().StringVar(&rktctx.gcGracePeriod, "grace-period", "", "Garbage grace period")
}
