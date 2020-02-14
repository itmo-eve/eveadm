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
	"log"

	"github.com/spf13/cobra"
)

// xenStopCmd represents the stop command
var xenStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Run shell command with arguments in 'stop' action on 'xen' mode",
	Long: `Run shell command with arguments in 'stop' action on 'xen' mode. For example:

eveadm xen stop uuid`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalf("Error in get param force in %s", cmd.Name())
		}
		xenctx.force = force
		xenctx.containerUUID = arg
		err, args, envs := xenctx.xenStopToCmd()
		if err != nil {
			log.Fatalf("Error in obtain params in %s", cmd.Name())
		}
		Run(cmd, Timeout, args, envs)
	},
}

func init() {
	xenCmd.AddCommand(xenStopCmd)
	xenStopCmd.Flags().BoolP("force", "f", false, "Force stop")
}
