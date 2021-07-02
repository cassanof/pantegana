package client

import (
	"os/exec"
	"runtime"
	"strings"
)

type SysInfo struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	Name string `json:"name"`

	// Linux only
	Distro string `json:"distro"`
	Kernel string `json:"kernel"`
}

var ClientSysInfo SysInfo

func (i *SysInfo) fingerprintLinux() {
	out, err := exec.Command("uname", "-rn").Output()
	if err == nil {
		split := strings.SplitN(string(out), " ", 2)
		i.Name = split[0]
		i.Kernel = split[1]
	}
}

// TODO: fingerprintWindows
func (i *SysInfo) fingerprintWindows() {
}

// TODO: fingerprintOsx
func (i *SysInfo) fingerprintOsx() {
}

func init() {
	ClientSysInfo = SysInfo{
		// run `go tool dist list` to show all options
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	switch runtime.GOOS {
	case "linux":
		go ClientSysInfo.fingerprintLinux()
	case "windows":
		go ClientSysInfo.fingerprintWindows()
	case "darwin":
		go ClientSysInfo.fingerprintOsx()
	default:
		// probably some kind of bsd so...
		go ClientSysInfo.fingerprintLinux()
	}
}
