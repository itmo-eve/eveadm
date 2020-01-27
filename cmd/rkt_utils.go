package cmd

type RKTContext struct {
	dir string
}

var rktctx RKTContext

func rktListToCmd(ctx RKTContext) (err error, args []string, envs string) {
	args = make([]string, 2)
	args[0] = "rkt"
	args[1] = "list"
	if ctx.dir != "" {
		args = append(args, "--dir="+ctx.dir)
	}
	envs = "a=1"
	err = nil
	return
}
