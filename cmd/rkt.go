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

	"github.com/spf13/cobra"
)

// rktCmd represents the rkt command
var rktCmd = &cobra.Command{
	Use:   "rkt",
	Short: "RKT mode",
	Long: `
Execute actions on 'rkt' mode. For example:

eveadm rkt list
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rkt called")
	},
}

func init() {
	rootCmd.AddCommand(rktCmd)
	rootCmd.PersistentFlags().StringVar(&rktctx.dir, "dir", "", "RKT data dir")
	rootCmd.PersistentFlags().StringVar(&rktctx.insecureOptions, "insecure-options", "image", "RKT insecure-options")
	rootCmd.PersistentFlags().StringVar(&rktctx.stage1Type, "stage1-type", "xen", "Type of stage1 (xen or general)")
}
