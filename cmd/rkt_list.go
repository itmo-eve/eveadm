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

// rktListCmd represents the list command
var rktListCmd = &cobra.Command{
	Use:   "list",
	Short: "Run shell command with arguments in 'list' action on 'rkt' mode",
	Long: `
Run shell command with arguments in 'list' action on 'rkt' mode. For example:

eveadm rkt list --image
eveadm rkt list
`,
	Run: func(cmd *cobra.Command, args []string) {
		isImage, err := cmd.Flags().GetBool("image")
		if err != nil {
			log.Fatalf("Error in get param image in %s", cmd.Name())
		}
		fields, err := cmd.Flags().GetString("fields")
		if err != nil {
			log.Fatalf("Error in get param fields in %s", cmd.Name())
		}
		rktctx.fields = fields
		var envs string
		if isImage {
			err, args, envs = rktctx.rktListImageToCmd()
		} else {
			err, args, envs = rktctx.rktListToCmd()
		}
		if err != nil {
			log.Fatalf("Error in obtain params in %s", cmd.Name())
		}
		Run(cmd, Timeout, args, envs)
	},
}

func init() {
	rktCmd.AddCommand(rktListCmd)
	rktListCmd.Flags().BoolVar(&rktctx.noLegend, "no-legend", false, "Suppress legend")
	rktListCmd.Flags().BoolP("image", "i", false, "Work with images")
	rktListCmd.Flags().String("fields", "", "Fields to return")
}
