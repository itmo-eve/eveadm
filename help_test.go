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
