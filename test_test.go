package main

import (
	"bytes"
//	"fmt"
	"strings"
	"testing"
	"os/exec"
	"time"

	"github.com/itmo-eve/eveadm/cmd"
)

func run_cmp (sys []string, eve []string, cmpf func(a, b string) int) int {
	c := exec.Command(sys[0], sys[1:]...)
	var out bytes.Buffer
	var err bytes.Buffer

	c.Stdout = &out
	c.Stderr = &err
	e := c.Run()
	cout := out.String()
	cerr := err.String()

	eout, _ := executeCommand(cmd.RootCmd, eve...)

	if e != nil {
		ret := strings.Compare(cerr, eout)
		return ret
	}

	ret := cmpf(cout, eout)
	return ret
}

func TestTestExecute (t *testing.T) {
	var eout string
	eout, _ = executeCommand(cmd.RootCmd, "test", "ls")
	checkStringContains(t, eout, "README.md")

	eout, _ = executeCommand(cmd.RootCmd, "test", "ls", "qwe")
	checkStringContains(t, eout, "No such file or directory")

	eout, _ = executeCommand(cmd.RootCmd, "test", "locale")
	checkStringOmits(t, eout, "LANG=ru_RU.UTF-8")
	//checkStringContains(t, eout, "LANG=C")

	eout, _ = executeCommand(cmd.RootCmd, "test", "locale", "-e", "LANG=ru_RU.UTF-8")
	checkStringContains(t, eout, "LANG=ru_RU.UTF-8")

	start := time.Now()
	eout, _ = executeCommand(cmd.RootCmd, "test", "sleep", "100")
	elapsed := time.Since(start)
	//er := elapsed.Round(100 * time.Second)
	if elapsed < 100 * time.Second {
		t.Errorf("Expected time of execution for 'sleep 100': \n %v\nGot:\n %v\n", 100 * time.Second, elapsed)
	}

	start = time.Now()
	eout, _ = executeCommand(cmd.RootCmd, "test", "sleep", "100", "-t", "1")
	elapsed = time.Since(start)
	//er = elapsed.Round(1 * time.Minute)
	if elapsed > 100 * time.Second {
		t.Errorf("Expected time of execution for 'sleep 100 -t 1': \n %v\nGot:\n %v\n", 1 * time.Minute, elapsed)
	}
}
