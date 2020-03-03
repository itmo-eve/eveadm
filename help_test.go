package eveadm

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/itmo-eve/eveadm/cmd"
)

var tests = map[string][]string{
	/*
	 */
	"tests/help":            {"help"},
	"tests/h":               {"-h"},
	"tests/help_rkt":        {"help", "rkt"},
	"tests/rkt-h":           {"rkt", "-h"},
	"tests/help_rkt_create": {"help", "rkt", "create"},
	"tests/rkt_create-h":    {"rkt", "create", "-h"},
	"tests/help_rkt_delete": {"help", "rkt", "delete"},
	"tests/rkt_delete-h":    {"rkt", "delete", "-h"},
	"tests/help_rkt_info":   {"help", "rkt", "info"},
	"tests/rkt_info-h":      {"rkt", "info", "-h"},
	"tests/help_rkt_list":   {"help", "rkt", "list"},
	"tests/rkt_list-h":      {"rkt", "list", "-h"},
	"tests/help_rkt_start":  {"help", "rkt", "start"},
	"tests/rkt_start-h":     {"rkt", "start", "-h"},
	"tests/help_rkt_stop":   {"help", "rkt", "stop"},
	"tests/rkt_stop-h":      {"rkt", "stop", "-h"},
	"tests/help_rkt_update": {"help", "rkt", "update"},
	"tests/rkt_update-h":    {"rkt", "update", "-h"},
	"tests/help_xen":        {"help", "xen"},
	"tests/xen-h":           {"xen", "-h"},
	"tests/help_xen_create": {"help", "xen", "create"},
	"tests/xen_create-h":    {"xen", "create", "-h"},
	"tests/help_xen_delete": {"help", "xen", "delete"},
	"tests/xen_delete-h":    {"xen", "delete", "-h"},
	"tests/help_xen_info":   {"help", "xen", "info"},
	"tests/xen_info-h":      {"xen", "info", "-h"},
	"tests/help_xen_list":   {"help", "xen", "list"},
	"tests/xen_list-h":      {"xen", "list", "-h"},
	"tests/help_xen_start":  {"help", "xen", "start"},
	"tests/xen_start-h":     {"xen", "start", "-h"},
	"tests/help_xen_stop":   {"help", "xen", "stop"},
	"tests/xen_stop-h":      {"xen", "stop", "-h"},
	"tests/help_xen_update": {"help", "xen", "update"},
	"tests/xen_update-h":    {"xen", "update", "-h"},
}

func TestHelpExecute(t *testing.T) {
	for f, a := range tests {
		name := strings.Join(a, " ")
		fmt.Println(name)
		t.Run(name, func(t *testing.T) {
			dat, err := ioutil.ReadFile(f)
			check(err)
			tst := string(dat)

			out, err := executeCommand(cmd.RootCmd, a...)
			check(err)
			res := strings.Compare(tst, out)
			if res != 0 {
				e := fmt.Sprintf("Command not passed -- args: %q file: %s\n", a, f)
				t.Errorf(e)
			}
		})
	}
}
