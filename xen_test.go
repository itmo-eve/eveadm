package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestXenSequence(t *testing.T) {
	var xenTestDir = "xen_pods_test"
	var binaryName = "eveadm"
	var testName = "testxen"
	var containerUUID = ""
	dir, err := os.Getwd()
	xlcmd := exec.Command("which", "xl")
	output, err := xlcmd.Output()
	if err != nil {
		t.Fatal("No xl found")
	}
	idcmd := exec.Command("id", "-u")
	output, err = idcmd.Output()

	if err != nil {
		t.Fatal(err)
	}
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		t.Fatal(err)
	}
	if i == 0 {
		dname, err := ioutil.TempDir("", xenTestDir)
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dname)
		t.Logf("Xen folder: %s", dname)
		out, err := os.Create(path.Join(dname, "cirros.qcow2"))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.Get("http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img")
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = out.Sync()
		if err != nil {
			t.Fatal(err)
		}
		err = resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = out.Close()
		if err != nil {
			t.Fatal(err)
		}
		configFile, err := os.Create(path.Join(dname, "config.cfg"))
		if err != nil {
			t.Fatal(err)
		}

		_, err = configFile.WriteString(`name = "` + testName + `"
on_poweroff = "preserve"
bootloader = "pygrub"
extra = "console=hvc0 root=/dev/xvda1"
memory = 128
vcpus = 1
vif = [ 'bridge=xenbr0' ]
disk = [ '` + path.Join(dname, "cirros.qcow2") + `,qcow2,xvda,rw' ]
`)
		if err != nil {
			t.Fatal(err)
		}
		err = configFile.Sync()
		if err != nil {
			t.Fatal(err)
		}
		err = configFile.Close()
		if err != nil {
			t.Fatal(err)
		}
		tests := []struct {
			name      string
			args      []string
			contains  []string
			omit      []string
			useEveAdm bool
			sleep     time.Duration
		}{
			{"container_create", []string{"xen", "create", "--xen-cfg-filename=" + path.Join(dname, "config.cfg"), "--paused"}, []string{}, []string{"error"}, true, 0},
			{"container_info_with_name", []string{"xen", "info", "--domname", testName}, []string{}, []string{"error"}, true, 10},
			{"container_start", []string{"xen", "start", containerUUID}, []string{}, []string{"error"}, true, 0},
			{"container_info", []string{"xen", "info", containerUUID}, []string{"domid"}, []string{"error"}, true, 0},
			{"container_stop", []string{"xen", "stop", containerUUID}, []string{}, []string{"error"}, true, 0},
			{"container_delete", []string{"xen", "delete", containerUUID}, []string{}, []string{"error"}, true, 10},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.sleep > 0 {
					time.Sleep(tt.sleep * time.Second)
				}

				cmd := exec.Command(path.Join(dir, binaryName), tt.args...)
				if !tt.useEveAdm {
					cmdFile := tt.args[0]
					cmdArgs := tt.args[1:]
					cmd = exec.Command(cmdFile, cmdArgs...)
				}
				t.Log(strings.Join(cmd.Args, " "))
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatal(err)
				}

				actual := strings.TrimSpace(string(output))

				for _, expected := range tt.contains {
					if !strings.Contains(actual, expected) {
						t.Fatalf("actual = %s, expected = %s", actual, expected)
					}
				}
				for _, omitted := range tt.omit {
					if strings.Contains(actual, omitted) {
						t.Fatalf("actual = %s, omitted = %s", actual, omitted)
					}
				}
				if tt.name == "container_info_with_name" {
					containerUUID = strings.TrimSpace(actual)
					//fix actual containerUUID
					tests[2].args = []string{"xen", "start", containerUUID}
					tests[3].args = []string{"xen", "info", containerUUID}
					tests[4].args = []string{"xen", "stop", containerUUID}
					tests[5].args = []string{"xen", "delete", containerUUID}
				}
			})
		}

	} else {
		t.Skipf("XEN test must be run as root! (sudo)")
	}
}
