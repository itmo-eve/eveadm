package cmd

import (
	"errors"
)

type XENContext struct {
	containerUUID  string
	xenCfgFilename string
	containerName  string
	runPaused      bool
	force          bool
}

var xenctx XENContext

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
	args = []string{"xl", "create", ctx.xenCfgFilename}
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
	args = []string{"xl", "domid", ctx.containerName}
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
