package main

import (
	"fmt"
	"context"
	"os/exec"
	"bytes"
	"time"

	"github.com/spf13/cobra"
)

var plugin_name = "xen"

func run(timeout time.Duration, args []string) {
	if len(args) > 0 {
		cmd := args[0]
		args = args[1:]
		if timeout != 0 {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			cmd := exec.CommandContext(ctx, cmd, args...)
			var sout bytes.Buffer
			var serr bytes.Buffer
			cmd.Stdout = &sout
			cmd.Stderr = &serr
			err := cmd.Run()
			if err != nil {
				fmt.Println("cmd: ", err.Error())
				cerr := ctx.Err()
				if cerr != nil {
					fmt.Println("ctx: ", cerr.Error())
				}
			}
			fmt.Printf("Output: %q\n", sout.String())
			fmt.Printf("Errors: %q\n", serr.String())
		} else {
			cmd := exec.Command(cmd, args...)
			var sout bytes.Buffer
			var serr bytes.Buffer
			cmd.Stdout = &sout
			cmd.Stderr = &serr
			err := cmd.Run()

			if err != nil {
				fmt.Println("cmd: ", err.Error())
			}
			fmt.Printf("Output: %q\n", sout.String())
			fmt.Printf("Errors: %q\n", serr.String())
		}
	}
}

func Console(timeout time.Duration, args []string) {
	fmt.Println(plugin_name, "'console' called", args)
	run(timeout, args)
}

func Create(timeout time.Duration, args []string) {
	fmt.Println(plugin_name, "'create' called", args)
	run(timeout, args)
}

func Delete(timeout time.Duration, args []string) {
	fmt.Println(plugin_name, "'delete' called", args)
	run(timeout, args)
}

func List(timeout time.Duration, args []string) {
	fmt.Println(plugin_name, "'list' called", args)
	run(timeout, args)
}

func Reboot(timeout time.Duration, args []string) {
	fmt.Println(plugin_name, "'reboot' called", args)
	run(timeout, args)
}

func Start(timeout time.Duration, args []string) {
	fmt.Println(plugin_name, "'start' called", args)
	run(timeout, args)
}

func Stop(timeout time.Duration, args []string) {
	fmt.Println(plugin_name, "'stop' called", args)
	run(timeout, args)
}

//func Init(cmd *cobra.Command) {
func Init(cmd *cobra.Command) {
	cmd.Short = fmt.Sprintf("Run shell command with arguments in '%s' action on '%s' mode", cmd.Name(), plugin_name)
	cmd.Use = fmt.Sprintln(cmd.Name(), "[shell_command args]")
	cmd.Long = fmt.Sprintf(`
Run shell command with arguments in '%s' action on '%s' mode. For example:

eveadm %s ps x
`, cmd.Name(), plugin_name, cmd.Name())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
