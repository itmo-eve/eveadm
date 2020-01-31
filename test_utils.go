package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(command *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(command, args...)
	return output, err
}

func executeCommandC(command *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	command.SetOutput(buf)
	command.SetArgs(args)

	c, err = command.ExecuteC()

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
