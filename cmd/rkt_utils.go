package cmd

import (
	"errors"
	"fmt"
	"strings"
)

type kvMap struct {
	mapping map[string]string
}

func (e *kvMap) Set(s string) error {
	if e.mapping == nil {
		e.mapping = make(map[string]string)
	}
	pair := strings.SplitN(s, "=", 2)
	if len(pair) != 2 {
		return fmt.Errorf("must be specified as key=value")
	}
	if _, exists := e.mapping[pair[0]]; exists {
		return fmt.Errorf("key %q already set", pair[0])
	}
	e.mapping[pair[0]] = pair[1]
	return nil
}

func (e *kvMap) IsEmpty() bool {
	return len(e.mapping) == 0
}

func (e *kvMap) String() string {
	return strings.Join(e.Strings(), "\n")
}

func (e *kvMap) Strings() []string {
	var env []string
	for n, v := range e.mapping {
		env = append(env, n+"="+v)
	}
	return env
}

func (e *kvMap) Type() string {
	return "kvMap"
}

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
	stage2MP        string
	flagExplicitEnv kvMap
	runPaused       bool
	stage1Path      string
	noOverlay       bool
	stage1Type      string
	force           bool
	gc              bool
	gcGracePeriod   string
	format          string
}

var rktctx RKTContext

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
	if ctx.uuidFile != "" {
		args = append(args, "--uuid-file-save="+ctx.uuidFile)
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
	if ctx.stage2MP != "" {
		envs += " STAGE2_MNT_PTS=" + ctx.stage2MP
	}
	explicitEnv := ctx.flagExplicitEnv.Strings()
	if len(explicitEnv) > 0 {
		for _, el := range explicitEnv {
			args = append(args, "--set-env="+el)
		}
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
func (ctx RKTContext) rktDeleteGC(isImage bool) (err error, args []string, envs string) {
	args = []string{"rkt"}
	if isImage {
		args = append(args, "image")
	}
	args = append(args, "gc")
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.gcGracePeriod != "" {
		args = append(args, "--grace-period="+ctx.gcGracePeriod)
	}
	envs = ""
	err = nil
	return
}
func (ctx RKTContext) rktDeleteGCImageToCmd() (err error, args []string, envs string) {
	return rktctx.rktDeleteGC(true)
}
func (ctx RKTContext) rktDeleteGCToCmd() (err error, args []string, envs string) {
	return rktctx.rktDeleteGC(false)
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
