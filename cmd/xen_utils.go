package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type XENContext struct {
	containerUUID  string
	xenCfgFilename string
	containerName  string
	runPaused      bool
	force          bool
}

var xenctx XENContext

func (ctx XENContext) xenRuneWrapper(timeout time.Duration, args []string, env string, cmdName string) {
	err, cerr, stdout, stderr := rune(timeout, args, env)
	if cerr != nil {
		log.Fatalf("Context error in %s", cmdName)
	}
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			fmt.Printf("%s", stdout.String())
			_, err = fmt.Fprintf(os.Stderr, "%s", stderr.String())
			if err != nil {
				fmt.Printf("%s", stderr.String())
			}
			os.Exit(waitStatus.ExitStatus())
		} else {
			_, err = fmt.Fprintf(os.Stderr, "Execute error in %s: %s\n", cmdName, err.Error())
			if err != nil {
				fmt.Printf("Execute error in %s: %s\n", cmdName, err.Error())
			}
		}
	}
	fmt.Printf("%s", stdout.String())
}

func (ctx XENContext) xenListToCmd() (err error, args []string, envs string) {
	args = []string{"xl", "list"}
	envs = ""
	err = nil
	return
}
func (ctx XENContext) xenCreateToCmd() (err error, args []string, envs string) {
	if ctx.xenCfgFilename == "" {
		return errors.New("No xenCfgFilename in args"), nil, ""
	}
	args = []string{"rkt", "run", ctx.xenCfgFilename}
	if ctx.runPaused {
		args = append(args, "-p")
	}
	err = nil
	return
}
func (ctx XENContext) xenStopToCmd() (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"xl", "shutdown"}
	if ctx.force {
		args = append(args, "-F")
	}
	args = append(args, ctx.containerUUID)
	envs = ""
	err = nil
	return
}
func (ctx XENContext) xenInfoToCmd() (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"xl", "list", "-l", ctx.containerUUID}
	envs = ""
	err = nil
	return
}
func (ctx XENContext) xenInfoDomidToCmd() (err error, args []string, envs string) {
	if ctx.containerName == "" {
		return errors.New("No container name in args"), nil, ""
	}
	args = []string{"xl", "domid", "-l", ctx.containerName}
	envs = ""
	err = nil
	return
}
func (ctx XENContext) xenDeleteToCmd() (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"xl", "destroy", ctx.containerUUID}
	envs = ""
	err = nil
	return
}
func (ctx XENContext) xenStartToCmd() (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"xl", "unpause", ctx.containerUUID}
	envs = ""
	err = nil
	return
}
