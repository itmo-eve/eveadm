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

type RKTContext struct {
	dir             string
	insecureOptions string
	noLegend        bool
	fields          string
	imageUUID       string
	containerUUID   string
	imageUrl        string
	uuidFile        string
	xenCfgFilename  string
	runPaused       bool
	stage1Path      string
	noOverlay       bool
	stage1Type      string
	force           bool
	format          string
}

var rktctx RKTContext

func (ctx RKTContext) rktRuneWrapper(timeout time.Duration, args []string, env string, cmdName string) {
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

func (ctx RKTContext) rktListToCmd() (err error, args []string, envs string) {
	args = []string{"rkt", "list"}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	if ctx.noLegend {
		args = append(args, "--no-legend")
	}
	envs = ""
	err = nil
	return
}

func (ctx RKTContext) rktListImageToCmd() (err error, args []string, envs string) {
	args = []string{"rkt", "image", "list"}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	if ctx.noLegend {
		args = append(args, "--no-legend")
	}
	if ctx.fields != "" {
		args = append(args, "--fields="+ctx.fields)
	}
	envs = ""
	err = nil
	return
}

func (ctx RKTContext) rktInfoImageToCmd() (err error, args []string, envs string) {
	if ctx.imageUUID == "" {
		return errors.New("No imageUUID in args"), nil, ""
	}
	args = []string{"rkt", "image", "cat-manifest"}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	args = append(args, ctx.imageUUID)
	envs = ""
	err = nil
	return
}
func (ctx RKTContext) rktCreateToCmd() (err error, args []string, envs string) {
	if ctx.imageUUID == "" {
		return errors.New("No image uuid in args"), nil, ""
	}
	args = []string{"rkt", "run", ctx.imageUUID}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	if ctx.stage1Path != "" {
		args = append(args, "--stage1-path="+ctx.stage1Path)
	}
	if ctx.noOverlay {
		args = append(args, "--no-overlay")
	}
	if ctx.runPaused {
		envs += " STAGE1_XL_OPTS=-p"
	}
	if ctx.xenCfgFilename != "" {
		envs += " STAGE1_SEED_XL_CFG=" + ctx.xenCfgFilename
	}
	err = nil
	return
}
func (ctx RKTContext) rktCreateImageToCmd() (err error, args []string, envs string) {
	if ctx.imageUrl == "" {
		return errors.New("No image url in args"), nil, ""
	}
	args = []string{"rkt", "fetch"}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	args = append(args, ctx.imageUrl)
	envs = ""
	err = nil
	return
}
func (ctx RKTContext) rktStopToCmd() (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"rkt", "stop", ctx.containerUUID}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	if ctx.force {
		args = append(args, "--force=true")
	}
	envs = ""
	err = nil
	return
}
func (ctx RKTContext) rktInfoToCmd() (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"rkt", "status", ctx.containerUUID}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	if ctx.format != "" {
		args = append(args, "--format="+ctx.format)
	}
	envs = ""
	err = nil
	return
}
func (ctx RKTContext) rktDeleteToCmd() (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"rkt", "rm", ctx.containerUUID}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	envs = ""
	err = nil
	return
}
func (ctx RKTContext) rktDeleteImageToCmd() (err error, args []string, envs string) {
	if ctx.imageUUID == "" {
		return errors.New("No image uuid in args"), nil, ""
	}
	args = []string{"rkt", "image", "rm", ctx.imageUUID}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	envs = ""
	err = nil
	return
}
func (ctx RKTContext) rktStartToCmd() (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"xl", "unpause", ctx.containerUUID}
	envs = ""
	err = nil
	return
}
