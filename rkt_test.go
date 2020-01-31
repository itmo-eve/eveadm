package main

import (
	rkt_cmd "github.com/itmo-eve/eveadm/cmd"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

var rktTestDir = "pods_test"

func TestRktSequence(t *testing.T) {
	idcmd := exec.Command("id", "-u")
	output, err := idcmd.Output()

	if err != nil {
		t.Fatal(err)
	}
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		t.Fatal(err)
	}
	var eout string
	if i == 0 {
		dname, err := ioutil.TempDir("", rktTestDir)
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dname)
		t.Logf("Rkt folder: %s", dname)
		eout, _ = executeCommand(rkt_cmd.RootCmd, "rkt", "create", "--image=true", "coreos.com/etcd:v2.0.0", "--dir=\""+dname+"\"", "--help=false")
		checkStringContains(t, eout, "sha512")
		if strings.Contains(eout, "sha512") {
			t.Logf("Rkt image sha: %s", eout)
		}
		imageHash := strings.TrimSpace(eout)
		eout, _ = executeCommand(rkt_cmd.RootCmd, "rkt", "list", "--image=true", "--dir=\""+dname+"\"", "--no-legend=true", "--help=false")
		checkStringContains(t, eout, "sha512")
		checkStringContains(t, eout, "coreos.com/etcd:v2.0.0")
		checkStringOmits(t, eout, "LAST USED")
		eout, _ = executeCommand(rkt_cmd.RootCmd, "rkt", "info", "--image=true", imageHash, "--dir=\""+dname+"\"", "--help=false")
		checkStringContains(t, eout, "coreos.com/etcd")
		checkStringContains(t, eout, "ImageManifest")
		t.Logf("info image: %s", eout)

	} else {
		t.Skipf("RKT test must be run as root! (sudo)")
	}
}
