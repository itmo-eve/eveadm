package cmd

import (
	"fmt"
	"context"
	"os"
	"os/exec"
	"bytes"
	"time"
	"strings"
)

var envs string

// Run shell command with arguments
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
				fmt.Println("Command error: ", err.Error())
				cerr := ctx.Err()
				if cerr != nil {
					fmt.Println("Command error:",
						cerr.Error())
				}
			}
			fmt.Printf("Stdout:\n%s\n", sout.String())
			fmt.Printf("Stderr:\n%s\n", serr.String())
		} else {
			cmd := exec.Command(cmd, args...)
			var sout bytes.Buffer
			var serr bytes.Buffer
			cmd.Stdout = &sout
			cmd.Stderr = &serr
			err := cmd.Run()

			if err != nil {
				fmt.Println("Command error:", err.Error())
			}
			fmt.Printf("Stdout:\n%s\n", sout.String())
			fmt.Printf("Stderr:\n%s\n", serr.String())
		}
	}
}

// Run shell command with arguments and enviroment variables
func rune(timeout time.Duration, args []string, env string) {
	envs := strings.Split(env, " ")
	if len(args) > 0 {
		cmd := args[0]
		args = args[1:]
		if timeout != 0 {
			ctx, cancel := context.WithTimeout(context.Background(),
				timeout)
			defer cancel()

			cmd := exec.CommandContext(ctx, cmd, args...)
			var sout bytes.Buffer
			var serr bytes.Buffer
			cmd.Stdout = &sout
			cmd.Stderr = &serr
			cmd.Env = append(os.Environ(), envs...)

			err := cmd.Run()
			if err != nil {
				fmt.Println("Command error: ", err.Error())
				cerr := ctx.Err()
				if cerr != nil {
					fmt.Println("Command error: ",
						cerr.Error())
				}
			}
			fmt.Printf("Stdout:\n%s\n", sout.String())
			fmt.Printf("Stderr:\n%s\n", serr.String())
		} else {
			cmd := exec.Command(cmd, args...)
			var sout bytes.Buffer
			var serr bytes.Buffer
			cmd.Stdout = &sout
			cmd.Stderr = &serr
			cmd.Env = append(os.Environ(), envs...)

			err := cmd.Run()

			if err != nil {
				fmt.Println("Command error: ", err.Error())
			}
			fmt.Printf("Stdout:\n%s\n", sout.String())
			fmt.Printf("Stderr:\n%s\n", serr.String())
		}
	}
}

