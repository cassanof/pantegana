package fingerprinter

import (
	"encoding/json"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var clientSysInfo SysInfo

func GetCurrentSysInfo() SysInfo {
	return clientSysInfo
}

// type declaration for the command functions
type command func(wg *sync.WaitGroup, i *SysInfo)

func Run() {
	clientSysInfo = SysInfo{
		// run `go tool dist list` to show all options
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	var wg sync.WaitGroup

	// loops through every command present for the OS (nix.go, darwin.go, win.go)
	// and runs the commands concurrently
	for _, cmd := range commands {
		wg.Add(1)
		go cmd(&wg, &clientSysInfo)
	}

	// Wait for commands to finish
	wg.Wait()

	dbg, _ := json.MarshalIndent(clientSysInfo, "", " ")
	log.Println(string(dbg))
}

// Helper func to execute commands, get output and handle errors
func runCmdAndGetOutput(expect int, cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	// if there is an error, it will return "unknown" as many times as the expect arg is defined as
	if err != nil {
		return "unknown" + strings.Repeat(" unknown\t", expect-1)
	}
	return trim(out)
}

// Helper func to trim '/n' and convert byte slice to string
func trim(out []byte) string {
	return strings.TrimSuffix(string(out), "\n")
}
