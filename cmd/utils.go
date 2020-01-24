package cmd

import (
	"fmt"
	"context"
	"os/exec"
	"bytes"
	"time"
)

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


