package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"testing"

	"github.com/itmo-eve/eveadm/cmd"
)

var tests = map[string][]string {
/*
*/
	"tests/help": []string {"help"},
        "tests/h": []string {"-h"},
        "tests/help_rkt": []string {"help", "rkt"},
        "tests/rkt-h": []string {"rkt", "-h"},
        "tests/help_rkt_create": []string {"help", "rkt", "create"},
        "tests/rkt_create-h": []string {"rkt", "create", "-h"},
        "tests/help_rkt_delete": []string {"help", "rkt", "delete"},
        "tests/rkt_delete-h": []string {"rkt", "delete", "-h"},
        "tests/help_rkt_info": []string {"help", "rkt", "info"},
        "tests/rkt_info-h": []string {"rkt", "info", "-h"},
        "tests/help_rkt_list": []string {"help", "rkt", "list"},
        "tests/rkt_list-h": []string {"rkt", "list", "-h"},
        "tests/help_rkt_start": []string {"help", "rkt", "start"},
        "tests/rkt_start-h": []string {"rkt", "start", "-h"},
        "tests/help_rkt_stop": []string {"help", "rkt", "stop"},
        "tests/rkt_stop-h": []string {"rkt", "stop", "-h"},
        "tests/help_rkt_update": []string {"help", "rkt", "update"},
        "tests/rkt_update-h": []string {"rkt", "update", "-h"},
        "tests/help_xen": []string {"help", "xen"},
        "tests/xen-h": []string {"xen", "-h"},
        "tests/help_xen_create": []string {"help", "xen", "create"},
        "tests/xen_create-h": []string {"xen", "create", "-h"},
        "tests/help_xen_delete": []string {"help", "xen", "delete"},
        "tests/xen_delete-h": []string {"xen", "delete", "-h"},
        "tests/help_xen_info": []string {"help", "xen", "info"},
        "tests/xen_info-h": []string {"xen", "info", "-h"},
        "tests/help_xen_list": []string {"help", "xen", "list"},
        "tests/xen_list-h": []string {"xen", "list", "-h"},
        "tests/help_xen_start": []string {"help", "xen", "start"},
        "tests/xen_start-h": []string {"xen", "start", "-h"},
        "tests/help_xen_stop": []string {"help", "xen", "stop"},
        "tests/xen_stop-h": []string {"xen", "stop", "-h"},
        "tests/help_xen_update": []string {"help", "xen", "update"},
        "tests/xen_update-h": []string {"xen", "update", "-h"},
}


func TestHelpExecute (t *testing.T) {
	for f, a := range tests {
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
	}
}
