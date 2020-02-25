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

// rktCreateCmd represents the create command
var rktCreateCmd = &cobra.Command{
	Use:   "create url/uuid",
	Short: "Run shell command with arguments in 'create' action on 'rkt' mode",
	Long: `
Run shell command with arguments in 'create' action on 'rkt' mode. For example:

eveadm rkt create --image url
eveadm rkt create image_uuid
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
			err, args, envs = rktctx.rktCreateImageToCmd()
		} else {
			rktctx.imageUUID = arg
			err, args, envs = rktctx.rktCreateToCmd()
		}
		if err != nil {
			log.Fatalf("Error in obtain params in %s", cmd.Name())
		}
		Run(cmd, Timeout, args, envs)
	},
}

func init() {
	rktCmd.AddCommand(rktCreateCmd)
	rktCreateCmd.Flags().BoolP("image", "i", false, "Work with images")
	rktCreateCmd.Flags().StringVar(&rktctx.uuidFile, "uuid-file-save", "", "File to save uuid")
	rktCreateCmd.Flags().StringVar(&rktctx.xenCfgFilename, "xen-cfg-filename", "", "File with xen cfg for stage1")
	rktCreateCmd.Flags().StringVar(&rktctx.stage1Path, "stage1-path", "/usr/sbin/stage1-xen.aci", "Stage1 path")
	rktCreateCmd.Flags().StringVar(&rktctx.stage2MP, "stage2-mnt-pts", "", "Stage2 mount points file")
	rktCreateCmd.Flags().Var(&rktctx.flagExplicitEnv, "set-env", "environment variable to set for all the apps in the form key=value")
	//Workaround to start in Ubuntu
	rktCreateCmd.Flags().BoolVar(&rktctx.noOverlay, "no-overlay", false, "Run without overlay")

	rktCreateCmd.Flags().BoolVar(&rktctx.runPaused, "paused", true, "Run paused")
}
