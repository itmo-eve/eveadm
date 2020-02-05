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

func TestRktSequence(t *testing.T) {

	type rktTestContext struct {
		containerHash string
		imageHash     string
	}
	type rktTestRun struct {
		name          string
		args          func() []string
		contains      func() []string
		omit          func() []string
		useEveAdm     bool
		sleep         time.Duration
		resultProcess func(result string)
	}

	var rktCtx rktTestContext

	var rktTestDir = "pods_test"
	var binaryName = "eveadm"
	var stage1XenPath = ""
	hasXL := false
	dir, err := os.Getwd()
	xlcmd := exec.Command("which", "xl")
	output, err := xlcmd.Output()
	if err != nil {
		t.Log("No xl found, run with default rkt stage1")
	} else {
		if strings.Contains(string(output), "xl") {
			hasXL = true
			stage1XenPath = path.Join(dir, "tests", "stage1-xen.aci")
			_, err := os.Stat(stage1XenPath)
			if os.IsNotExist(err) {
				t.Fatal("No stage1-xen.aci found in tests directory")
			}
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

		var tests []rktTestRun
		imageCreateTest := rktTestRun{
			"image_create",
			func() []string {
				return []string{"rkt", "create", "--image=true", "docker://library/alpine:3.11", "--dir=" + dname}
			},
			func() []string {
				return []string{"sha512", rktCtx.imageHash}
			}, func() []string {
				return []string{}
			},
			true,
			0,
			func(result string) { rktCtx.imageHash = result }}
		tests = append(tests, imageCreateTest)
		imageInfo := rktTestRun{
			"image_info",
			func() []string {
				return []string{"rkt", "info", "--image=true", rktCtx.imageHash, "--dir=" + dname}
			},
			func() []string {
				return []string{}
			},
			func() []string {
				return []string{}
			},
			true,
			0,
			func(result string) {}}
		tests = append(tests, imageInfo)
		containerCreate := rktTestRun{
			"container_create",
			func() []string {
				return []string{"systemd-run", path.Join(dir, binaryName), "rkt", "create", rktCtx.imageHash, "--dir=" + dname, "--no-overlay=true", "--stage1-path="}
			},
			func() []string {
				return []string{"Running as unit"}
			},
			func() []string {
				return []string{}
			},
			false,
			0,
			func(result string) {}}
		if hasXL {
			containerCreate = rktTestRun{
				"container_create",
				func() []string {
					return []string{"rkt", "create", rktCtx.imageHash, "--dir=" + dname, "--paused=true", "--stage1-path=" + stage1XenPath}
				},
				func() []string {
					return []string{}
				}, func() []string {
					return []string{"Running as unit"}
				},
				true,
				0,
				func(result string) {}}
		}
		tests = append(tests, containerCreate)
		containerList := rktTestRun{
			"container_list",
			func() []string {
				return []string{"rkt", "list", "--dir=" + dname, "--no-legend=true"}
			},
			func() []string {
				return []string{"alpine"}
			},
			func() []string {
				return []string{"ID"}
			},
			true,
			10,
			func(result string) {
				rktCtx.containerHash = strings.Split(result, "\t")[0]
			}}
		tests = append(tests, containerList)
		containerInfo := rktTestRun{
			"container_info",
			func() []string {
				return []string{"rkt", "info", "--dir=" + dname, rktCtx.containerHash}
			},
			func() []string {
				return []string{"state", "alpine"}
			},
			func() []string {
				return []string{}
			},
			true,
			0,
			func(result string) {
				re := regexp.MustCompile(`"name":[ ]?"([a-z0-9-]*)"`)
				resultsFind := re.FindSubmatch([]byte(result))
				if resultsFind != nil && len(resultsFind) == 2 {
					rktCtx.containerHash = string(resultsFind[1])
				}
			}}
		tests = append(tests, containerInfo)
		containerStart := rktTestRun{
			"container_start",
			func() []string {
				return []string{"rkt", "start", "--dir=" + dname, "--stage1-type=common", rktCtx.containerHash}
			},
			func() []string {
				return []string{"Not implemented"}
			},
			func() []string {
				return []string{}
			},
			true,
			0,
			func(result string) {}}
		if hasXL {
			containerStart = rktTestRun{
				"container_start",
				func() []string {
					return []string{"rkt", "start", "--dir=" + dname, "--stage1-type=xen", rktCtx.containerHash}
				},
				func() []string {
					return []string{}
				},
				func() []string {
					return []string{"Not implemented"}
				},
				true,
				0,
				func(result string) {}}
		}
		tests = append(tests, containerStart)
		containerStop := rktTestRun{
			"container_stop",
			func() []string {
				return []string{"rkt", "stop", "--dir=" + dname, rktCtx.containerHash}
			},
			func() []string {
				return []string{rktCtx.containerHash}
			},
			func() []string {
				return []string{}
			},
			true,
			0,
			func(result string) {}}
		tests = append(tests, containerStop)
		containerDelete := rktTestRun{
			"container_delete",
			func() []string {
				return []string{"rkt", "delete", "--dir=" + dname, rktCtx.containerHash}
			},
			func() []string {
				return []string{rktCtx.containerHash}
			},
			func() []string {
				return []string{}
			},
			true,
			0,
			func(result string) {}}
		tests = append(tests, containerDelete)
		imageDelete := rktTestRun{
			"image_delete",
			func() []string {
				return []string{"rkt", "delete", "--image=true", "--dir=" + dname, rktCtx.imageHash}
			},
			func() []string {
				return []string{rktCtx.imageHash}
			},
			func() []string {
				return []string{}
			},
			true,
			0,
			func(result string) {}}
		tests = append(tests, imageDelete)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.sleep > 0 {
					time.Sleep(tt.sleep * time.Second)
				}

				args := tt.args()
				cmd := exec.Command(path.Join(dir, binaryName), args...)
				if !tt.useEveAdm {
					cmdFile := args[0]
					cmdArgs := args[1:]
					cmd = exec.Command(cmdFile, cmdArgs...)
				}
				t.Log(strings.Join(cmd.Args, " "))
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatal(err)
				}

				actual := strings.TrimSpace(string(output))

				for _, expected := range tt.contains() {
					if !strings.Contains(actual, expected) {
						t.Fatalf("actual = %s, expected = %s", actual, expected)
					}
				}
				for _, omitted := range tt.omit() {
					if strings.Contains(actual, omitted) {
						t.Fatalf("actual = %s, omitted = %s", actual, omitted)
					}
				}
				t.Log(actual)
				tt.resultProcess(actual)
			})
		}

	} else {
		t.Skipf("RKT test must be run as root! (sudo)")
	}
}
