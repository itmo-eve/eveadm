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

func run_out (rerr error, cerr error, sout bytes.Buffer, serr bytes.Buffer) {
	if rerr != nil {
		fmt.Println("Command error: ", rerr.Error())
		if cerr != nil {
			fmt.Println("Command error: ",
				cerr.Error())
		}
	}
	fmt.Printf("Stdout:\n%s\n", sout.String())
	fmt.Printf("Stderr:\n%s\n", serr.String())
}

// Run shell command with arguments and enviroment variables
func rune(timeout time.Duration, args []string, env string) (rerr error, cerr error, stdout bytes.Buffer, stderr bytes.Buffer) {
	var re error
	var ce error
	var sout bytes.Buffer
	var serr bytes.Buffer

	if len(args) > 0 {
		envs := strings.Split(env, " ")
		cmd := args[0]
		args = args[1:]
		if timeout != 0 {
			var cancel context.CancelFunc
			cnt, cancel := context.WithTimeout(context.Background(),
				timeout)
			defer cancel()

			cmd := exec.CommandContext(cnt, cmd, args...)
			cmd.Stdout = &sout
			cmd.Stderr = &serr
			cmd.Env = append(os.Environ(), envs...)

			re = cmd.Run()
			ce = cnt.Err()
		} else {
			cmd := exec.Command(cmd, args...)
			cmd.Stdout = &sout
			cmd.Stderr = &serr
			cmd.Env = append(os.Environ(), envs...)

			re = cmd.Run()
		}
		if verbose {
			run_out(re, ce, sout, serr)
		}
	}
	return re, ce, sout, serr
}
