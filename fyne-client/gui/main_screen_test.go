package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "github.com/mitchellh/mapstructure"
	"github.com/zcalusic/sysinfo"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"testing"

	"strings"
)

func TestOsToTwoDimArray(t *testing.T) {
	fmt.Println("arch", runtime.GOARCH)
	current, err := user.Current()
	handleError(err)
	fmt.Println("user", current)
	fmt.Println("user", current.Uid)
	var si sysinfo.SysInfo

	data, err := json.MarshalIndent(&si, "", "  ")
	handleError(err)
	fmt.Println(string(data))

	env := os.Environ()
	fmt.Println("found his many env vars ", len(env))
	// for i,v :=range env {
	// 	fmt.Println("env var ", i,v)
	// }
}

func TestOsInfo(t *testing.T) {

	cmd := exec.Command("uname", "-srio")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	handleError(err)
	//fmt.Println(out.String())

	for i, v := range strings.Split(out.String(), " ") {
		fmt.Println("", i, v)
	}

}
