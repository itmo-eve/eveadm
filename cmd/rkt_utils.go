package cmd

type RKTContext struct {
	dir             string
	insecureOptions string
	noLegend        bool
	fields          string
	uuid            string
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
	args = []string{"rkt", "image", "cat-manifest"}
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	if ctx.insecureOptions != "" {
		args = append(args, "--insecure-options="+ctx.insecureOptions)
	}
	args = append(args, ctx.uuid)
	envs = ""
	err = nil
	return
}
