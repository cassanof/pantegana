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
	go func() {
		out, err := exec.Command("uname", "-rn").Output()
		if err == nil {
			split := strings.SplitN(trim(out), " ", 2)
			i.Name = split[0]
			i.Kernel = split[1]
		}
		defer wg.Done()
	}()
	go func() {
		out, err := exec.Command("lsb_release", "-d").Output()
		if err == nil {
			i.Distro = strings.SplitN(trim(out), "\t", 2)[1]
		}
		defer wg.Done()
	}()
	go func() {
		out, err := exec.Command("whoami").Output()
		if err == nil {
			i.User.Name = trim(out)
		}
		defer wg.Done()
	}()
	go func() {
		out, err := exec.Command("id", "-u").Output()
		if err == nil {
			i.User.Id = trim(out)
		}
		defer wg.Done()
	}()
	go func() {
		out, err := exec.Command("groups").Output()
		if err == nil {
			i.User.Groups = trim(out)
		}
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
		// probably some kind of bsd so...
		wg.Add(5) // Set here the number of commands to execute
		go clientSysInfo.fingerprintLinux(&wg)
	}

	// Wait for commands to finish
	wg.Wait()

	dbg, _ := json.MarshalIndent(clientSysInfo, "", " ")
	log.Println(string(dbg))
}

// Helper func to trim '/n' and convert byte arr to string
func trim(out []byte) string {
	return strings.TrimSuffix(string(out), "\n")
}
