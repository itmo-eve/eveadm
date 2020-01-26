package cmd

func rktListToCmd() (err error, args []string, envs string) {
	envs = ""
	err = nil
	args = make([]string, 2)
	args[0] = "rkt"
	args[1] = "list"
	return
}
