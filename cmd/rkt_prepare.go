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

// rktPrepareCmd represents the prepare command
var rktPrepareCmd = &cobra.Command{
	Use:   "prepare uuid",
	Short: "Run shell command with arguments in 'prepare' action on 'rkt' mode",
	Long: `
Run shell command with arguments in 'prepare' action on 'rkt' mode. For example:

eveadm rkt prepare uuid
`, Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uuid := args[0]
		rktctx.imageUUID = uuid
		err, args, envs := rktctx.rktPrepareImageToCmd()
		if err != nil {
			log.Fatalf("Error in obtain params in %s", cmd.Name())
		}
		Run(cmd, Timeout, args, envs)
	},
}

func init() {
	rktCmd.AddCommand(rktPrepareCmd)
	rktPrepareCmd.Flags().StringVar(&rktctx.format, "format", "json", "Format of output")
	rktPrepareCmd.Flags().BoolVar(&rktctx.noOverlay, "no-overlay", false, "Run without overlay")
	rktPrepareCmd.Flags().BoolVar(&rktctx.quiet, "quiet", false, "Run in quiet mode")
	rktPrepareCmd.Flags().StringVar(&rktctx.prepareName, "name", "", "Name for prepare")
	rktPrepareCmd.Flags().StringVar(&rktctx.stage1Path, "stage1-path", "/usr/sbin/stage1-xen.aci", "Stage1 path")
}
