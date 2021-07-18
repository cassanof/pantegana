package client

import (
	"encoding/json"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

type SysInfo struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	Name string `json:"name"`

	User struct {
		Name string `json:"name"`
		Id   string `json:"uid,omitempty"` // *nix only
	}

	// *nix only
	Distro string `json:"distro,omitempty"`
	Kernel string `json:"kernel,omitempty"`
}

var clientSysInfo SysInfo

func GetCurrentSysInfo() SysInfo {
	return clientSysInfo
}

func (i *SysInfo) fingerprintLinux() {
	out, err := exec.Command("uname", "-rn").Output()
	if err == nil {
		split := strings.SplitN(trim(out), " ", 2)
		i.Name = split[0]
		i.Kernel = split[1]
	}
	out, err = exec.Command("whoami").Output()
	if err == nil {
		i.User.Name = trim(out)
	}
	out, err = exec.Command("id", "-u").Output()
	if err == nil {
		i.User.Id = trim(out)
	}

	dbg, _ := json.MarshalIndent(i, "", " ")
	log.Println(string(dbg))
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

	switch runtime.GOOS {
	case "linux":
		go clientSysInfo.fingerprintLinux()
	case "windows":
		go clientSysInfo.fingerprintWindows()
	case "darwin":
		go clientSysInfo.fingerprintOsx()
	default:
		// probably some kind of bsd so...
		go clientSysInfo.fingerprintLinux()
	}
}

// Helper func to trim '/n' and convert byte arr to string
func trim(out []byte) string {
	return strings.TrimSuffix(string(out), "\n")
}
