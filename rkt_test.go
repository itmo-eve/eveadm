package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

var rktTestDir = "pods_test"
var binaryName = "eveadm"
var imageHash = "sha512-5bee98b12bc8eb63425d8af93a8dc1be"

func TestRktSequence(t *testing.T) {
	hasXL := false
	dir, err := os.Getwd()
	xlcmd := exec.Command("which", "xl")
	output, err := xlcmd.Output()
	if err != nil {
		t.Log("No xl found, run with default rkt stage1")
	} else {
		if strings.Contains(string(output), "xl") {
			hasXL = true
		}
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
		dname, err := ioutil.TempDir("", rktTestDir)
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dname)
		t.Logf("Rkt folder: %s", dname)

		containerHash := ""

		tests := []struct {
			name      string
			args      []string
			contains  []string
			omit      []string
			useEveAdm bool
			sleep     time.Duration
		}{
			{"image_create", []string{"rkt", "create", "--image=true", "docker://library/alpine:3.11", "--dir=" + dname}, []string{"sha512", imageHash}, []string{}, true, 0},
			{"image_list", []string{"rkt", "list", "--image=true", "--dir=" + dname, "--no-legend=true"}, []string{"sha512", "registry-1.docker.io/library/alpine:3.11"}, []string{"LAST USED"}, true, 0},
			{"image_info", []string{"rkt", "info", "--image=true", imageHash, "--dir=" + dname}, []string{}, []string{}, true, 0},
			{"container_create", []string{"systemd-run", path.Join(dir, binaryName), "rkt", "create", imageHash, "--dir=" + dname, "--no-overlay=true", "--stage1-path="}, []string{"Running as unit"}, []string{}, false, 0},
			{"container_list", []string{"rkt", "list", "--dir=" + dname, "--no-legend=true"}, []string{"alpine"}, []string{"ID"}, true, 10},
			{"container_info", []string{"rkt", "info", "--dir=" + dname, containerHash}, []string{"state"}, []string{}, true, 0},
			{"container_start", []string{"rkt", "start", "--dir=" + dname, "--stage1-type=common", containerHash}, []string{"Not implemented"}, []string{}, true, 0},
			{"container_stop", []string{"rkt", "stop", "--dir=" + dname, containerHash}, []string{containerHash}, []string{}, true, 0},
			{"container_delete", []string{"rkt", "delete", "--dir=" + dname, containerHash}, []string{containerHash}, []string{}, true, 0},
			{"image_delete", []string{"rkt", "delete", "--image=true", "--dir=" + dname, imageHash}, []string{imageHash}, []string{}, true, 0},
		}
		if hasXL {
			//fix container_create for stage1-xen.aci
			tests[3].args = []string{"rkt", "create", imageHash, "--dir=" + dname, "--paused=true"}
			tests[3].useEveAdm = true
			tests[3].contains = []string{}
			tests[3].omit = []string{"Running as unit"}
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
				if tt.name == "container_list" {
					containerHash = strings.Split(actual, "\t")[0]
					//fix container_start for actual containerHash
					tests[5].args = []string{"rkt", "info", "--dir=" + dname, containerHash}
				}
				if tt.name == "container_info" {
					re := regexp.MustCompile(`"name":[ ]?"([a-z0-9-]*)"`)
					resultsFind := re.FindSubmatch([]byte(actual))
					if resultsFind != nil && len(resultsFind) == 2 {
						containerHash = string(resultsFind[1])
						if hasXL {
							tests[6].args = []string{"rkt", "start", "--dir=" + dname, containerHash}
							tests[6].contains = []string{}
							tests[6].omit = []string{"Not implemented"}
						} else {
							tests[6].args = []string{"rkt", "start", "--dir=" + dname, "--stage1-type=common", containerHash}
						}
						tests[7].args = []string{"rkt", "stop", "--dir=" + dname, containerHash}
						tests[7].contains = []string{containerHash}
						tests[8].args = []string{"rkt", "delete", "--dir=" + dname, containerHash}
						tests[8].contains = []string{containerHash}
					}
				}
			})
		}

	} else {
		t.Skipf("RKT test must be run as root! (sudo)")
	}
}
