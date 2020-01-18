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
	"os"
	"plugin"
	"time"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type ActionFunc func(timeout time.Duration, args []string)

var Functions = map[string]ActionFunc {
	"Console":nil,
	"Create":nil,
	"Delete":nil,
	"List":nil,
	"Reboot":nil,
	"Start":nil,
	"Stop":nil,
}
// Cobra-module related functions
var Init func(consoleCmd *cobra.Command)

var cfgFile string
var Plugin string
var Plugins_dir string
var Timeout time.Duration
var timeout string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eveadm",
	Short: "Manage EVE virtual machines and containers",
	Long: `The eveadm tool allows you to interact with virtual machines
and containers on a EVE system. It allows you to create, inspect, modify and
delete virtual machines on the local system.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside %s PersistentPreRun %s with args: %v\n", cmd.Name(), args)
		if Init != nil {
			Init(consoleCmd)
			Init(createCmd)
			Init(deleteCmd)
			Init(listCmd)
			Init(rebootCmd)
			Init(startCmd)
			Init(stopCmd)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.eveadm.yaml)")
	rootCmd.PersistentFlags().StringVarP(&timeout, "timeout", "t", "", "Actions timeout in minutes")
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	rootCmd.PersistentFlags().StringVarP(&Plugins_dir, "plugins_dir", "", "plugins", "Plugin modules directory")
	viper.BindPFlag("plugins_dir", rootCmd.PersistentFlags().Lookup("plugins_dir"))
	rootCmd.PersistentFlags().StringVarP(&Plugin, "plugin", "p", "", "Plugin module")
	viper.BindPFlag("plugin", rootCmd.PersistentFlags().Lookup("plugin"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".eveadm" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".eveadm")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	timeout = viper.GetString("timeout")
	if len(timeout) > 0 {
		minutes, err := strconv.Atoi(timeout)
		Timeout = time.Duration(minutes) * time.Minute
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Timeout:", Timeout)

	Plugins_dir = viper.GetString("plugins_dir")
	Plugin = viper.GetString("plugin")
	name := fmt.Sprintf("%s/%s.so", Plugins_dir, Plugin)
	fmt.Println("Plugin:", name)

	p, err := plugin.Open(name)
	if err != nil {
		fmt.Println("Can't open", name)
		fmt.Println(err)
		os.Exit(1)
	}

	fs, err := p.Lookup("Init")
	if err == nil {
		f, ok := fs.(func(consoleCmd *cobra.Command))
		if ok {
			Init = f
		} else {
			fmt.Println(ok)
		}
	}

	for fn, _ := range Functions {
		fs, err := p.Lookup(fn)
		if err == nil {
			f, ok := fs.(func(timeout time.Duration, args []string))
			if ok {
				Functions[fn] = f
			}
		}
	}
}
