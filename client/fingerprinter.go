package client

import (
	"encoding/json"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type SysInfo struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	Name string `json:"name"`

	User struct {
		Name   string `json:"name"`
		Id     string `json:"uid,omitempty"`    // *nix only
		Groups string `json:"groups,omitempty"` // *nix only
	}

	// *nix only
	Distro string `json:"distro,omitempty"`
	Kernel string `json:"kernel,omitempty"`
}

var clientSysInfo SysInfo

func GetCurrentSysInfo() SysInfo {
	return clientSysInfo
}

func (i *SysInfo) fingerprintLinux(wg *sync.WaitGroup) {
	wg.Add(5) // Set here the number of commands to execute
	go func() {
		outs := strings.SplitN(runCmdAndGetOutput(2, "uname", "-rn"), " ", 2)
		i.Name = outs[0]
		i.Kernel = outs[1]
		defer wg.Done()
	}()
	go func() {
		i.Distro = strings.SplitN(runCmdAndGetOutput(3, "lsb_release", "-d"), "\t", 2)[1]
		defer wg.Done()
	}()
	go func() {
		i.User.Name = runCmdAndGetOutput(1, "whoami")
		defer wg.Done()
	}()
	go func() {
		i.User.Id = runCmdAndGetOutput(1, "id", "-u")
		defer wg.Done()
	}()
	go func() {
		i.User.Groups = runCmdAndGetOutput(1, "groups")
		defer wg.Done()
	}()
}

// TODO: fingerprintWindows
func (i *SysInfo) fingerprintWindows() {
}

// TODO: fingerprintOsx
func (i *SysInfo) fingerprintOsx() {
}

func RunFingerprinter() {
	clientSysInfo = SysInfo{
		// run `go tool dist list` to show all options
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	var wg sync.WaitGroup

	switch runtime.GOOS {
	case "windows":
		go clientSysInfo.fingerprintWindows()
	case "darwin":
		go clientSysInfo.fingerprintOsx()
	default:
		// probably some kind of *nix so...
		clientSysInfo.fingerprintLinux(&wg)
	}

	// Wait for commands to finish
	wg.Wait()

	dbg, _ := json.MarshalIndent(clientSysInfo, "", " ")
	log.Println(string(dbg))
}

// Helper func to execute commadns, get output and handle errors
func runCmdAndGetOutput(expect int, cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	// if there is an error, it will return "unknown" as many times as the expect arg is defined as
	if err != nil {
		return "unknown" + strings.Repeat(" unknown\t", expect-1)
	}
	return trim(out)
}

// Helper func to trim '/n' and convert byte arr to string
func trim(out []byte) string {
	return strings.TrimSuffix(string(out), "\n")
}
