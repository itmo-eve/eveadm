package main

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	//"github.com/stretchr/testify/assert"
	"github.com/itmo-eve/eveadm/cmd"
)

var tests = map[string][]string {
	"tests/help": []string {"help"},
        "tests/h": []string {"-h"},
        "tests/help_test": []string {"help", "test"},
        "tests/test-h": []string {"test", "-h"},
	"tests/help_test_create": []string {"help", "test", "create"},
        "tests/test_create-h": []string {"test", "create", "-h"},
        "tests/help_test_delete": []string {"help", "test", "delete"},
        "tests/test_delete-h": []string {"test", "delete", "-h"},
        "tests/help_test_info": []string {"help", "test", "info"},
        "tests/test_info-h": []string {"test", "info", "-h"},
        "tests/help_test_list": []string {"help", "test", "list"},
        "tests/test_list-h": []string {"test", "list", "-h"},
        "tests/help_test_start": []string {"help", "test", "start"},
        "tests/test_start-h": []string {"test", "start", "-h"},
        "tests/help_test_stop": []string {"help", "test", "stop"},
        "tests/test_stop-h": []string {"test", "stop", "-h"},
        "tests/help_test_update": []string {"help", "test", "update"},
        "tests/test_update-h": []string {"test", "update", "-h"},
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

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func checkStringContains(t *testing.T, got, expected string) {
	if !strings.Contains(got, expected) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	}
}

func checkStringOmits(t *testing.T, got, expected string) {
	if strings.Contains(got, expected) {
		t.Errorf("Expected to not contain: \n %v\nGot: %v", expected, got)
	}
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func TestExecute (t *testing.T) {
	for f, a := range tests {
		fmt.Printf("args: %q file: %s\n", f, a)
		dat, err := ioutil.ReadFile(f)
		check(err)
		tst := string(dat)

		out, err := executeCommand(cmd.RootCmd, a...)
		check(err)
		res := strings.Compare(tst, out)
		fmt.Println("res:", res)
		if res != 0 {
			t.Errorf("Command 'help' not passed")	
		}
	}
}
