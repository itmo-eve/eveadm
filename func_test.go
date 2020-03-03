package eveadm

import (
	//	"io/ioutil"
	"fmt"
	//	"strings"
	"testing"
	"time"

	"github.com/giggsoff/eveadm/cmd"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run shell command with arguments on 'test' mode",
	Long: `
Run shell command with arguments on 'test' mode. For example:

eveadm test ps x`,
	Run: func(c *cobra.Command, args []string) {
		cmd.Envs = c.Flag("env").Value.String()
		if cmd.Verbose {
			fmt.Println("test called with envs:", cmd.Envs)
		}
		cmd.Test = true
		cmd.Run(c, cmd.Timeout, args, cmd.Envs)
	},
}

func init() {
	cmd.RootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().StringP("env", "e", "", "Setting environment variables")
}

func TestFuncExecute(t *testing.T) {
	var eout string

	t.Run("correct execution", func(t *testing.T) {
		eout, _ = executeCommand(cmd.RootCmd, "test", "ls")
		checkStringContains(t, eout, "README.md")
	})

	t.Run("wrong execution", func(t *testing.T) {
		eout, _ = executeCommand(cmd.RootCmd, "test", "ls", "qwe")
		checkStringContains(t, eout, "No such file or directory")

	})

	t.Run("environment variables setting", func(t *testing.T) {
		eout, _ = executeCommand(cmd.RootCmd, "test", "locale")
		checkStringOmits(t, eout, "LANG=ru_RU.UTF-8")
		//checkStringContains(t, eout, "LANG=C")

		eout, _ = executeCommand(cmd.RootCmd, "test", "locale", "-e",
			"LANG=ru_RU.UTF-8")
		checkStringContains(t, eout, "LANG=ru_RU.UTF-8")

	})

	t.Run("timewait", func(t *testing.T) {
		start := time.Now()
		eout, _ = executeCommand(cmd.RootCmd, "test", "sleep", "100")
		elapsed := time.Since(start)
		if elapsed < 100*time.Second {
			t.Errorf("Expected time of execution for 'sleep 100': \n %v\nGot:\n %v\n", 100*time.Second, elapsed)
		}

		start = time.Now()
		eout, _ = executeCommand(cmd.RootCmd, "test", "sleep", "100", "-t", "1")
		elapsed = time.Since(start)
		if elapsed > 100*time.Second {
			t.Errorf("Expected time of execution for 'sleep 100 -t 1': \n %v\nGot:\n %v\n", 1*time.Minute, elapsed)
		}
	})
}
