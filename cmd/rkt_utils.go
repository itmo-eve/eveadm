package cmd

import "errors"

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
}

var rktctx RKTContext

func rktListToCmd(ctx RKTContext) (err error, args []string, envs string) {
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

func rktListImageToCmd(ctx RKTContext) (err error, args []string, envs string) {
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

func rktInfoImageToCmd(ctx RKTContext) (err error, args []string, envs string) {
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
func rktCreateToCmd(ctx RKTContext) (err error, args []string, envs string) {
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
func rktCreateImageToCmd(ctx RKTContext) (err error, args []string, envs string) {
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
func rktStopToCmd(ctx RKTContext) (err error, args []string, envs string) {
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
	envs = ""
	err = nil
	return
}
func rktInfoToCmd(ctx RKTContext) (err error, args []string, envs string) {
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
	envs = ""
	err = nil
	return
}
func rktDeleteToCmd(ctx RKTContext) (err error, args []string, envs string) {
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
func rktDeleteImageToCmd(ctx RKTContext) (err error, args []string, envs string) {
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
func rktStartToCmd(ctx RKTContext) (err error, args []string, envs string) {
	if ctx.containerUUID == "" {
		return errors.New("No container uuid in args"), nil, ""
	}
	args = []string{"xl", "unpause", ctx.containerUUID}
	envs = ""
	err = nil
	return
}
