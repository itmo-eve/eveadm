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

// rktInfoCmd represents the info command
var rktInfoCmd = &cobra.Command{
	Use:   "info uuid",
	Short: "Run shell command with arguments in 'info' action on 'rkt' mode",
	Long: `
Run shell command with arguments in 'info' action on 'rkt' mode. For example:

eveadm rkt info --image uuid
eveadm rkt info uuid
`, Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uuid := args[0]
		isImage, err := cmd.Flags().GetBool("image")
		if err != nil {
			log.Fatalf("Error in get param image in %s", cmd.Name())
		}
		var envs string
		if isImage {
			rktctx.imageUUID = uuid
			err, args, envs = rktctx.rktInfoImageToCmd()
		} else {
			rktctx.containerUUID = uuid
			err, args, envs = rktctx.rktInfoToCmd()
		}
		if err != nil {
			log.Fatalf("Error in obtain params in %s", cmd.Name())
		}
		Run(cmd, Timeout, args, envs)
	},
}

func init() {
	rktCmd.AddCommand(rktInfoCmd)
	rktInfoCmd.Flags().BoolP("image", "i", false, "Work with images")
	rktInfoCmd.Flags().StringVar(&rktctx.format, "format", "json", "Format of output")
}
