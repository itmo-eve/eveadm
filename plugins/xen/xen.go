package main

import (
	"fmt"
	"context"
	"os/exec"
	"bytes"
	"time"
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
