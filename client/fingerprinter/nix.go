//go:build linux
// +build linux

package fingerprinter

import (
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
	} `json:"user"`

	// *nix only
	Distro string `json:"distro,omitempty"`
	Kernel string `json:"kernel,omitempty"`
}

// holds all of the command functions to run
var commands []command = []command{
	uname,
	distro,
	whoami,
	userId,
	groups,
}

func uname(wg *sync.WaitGroup, i *SysInfo) {
	outs := strings.SplitN(runCmdAndGetOutput(2, "uname", "-rn"), " ", 2)
	i.Name = outs[0]
	i.Kernel = outs[1]
	defer wg.Done()
}

func distro(wg *sync.WaitGroup, i *SysInfo) {
	outs := strings.SplitN(runCmdAndGetOutput(2, "uname", "-rn"), " ", 2)
	i.Name = outs[0]
	i.Kernel = outs[1]
	defer wg.Done()
}

func whoami(wg *sync.WaitGroup, i *SysInfo) {
	i.User.Name = runCmdAndGetOutput(1, "whoami")
	defer wg.Done()
}

func userId(wg *sync.WaitGroup, i *SysInfo) {
	i.User.Id = runCmdAndGetOutput(1, "id", "-u")
	defer wg.Done()
}

func groups(wg *sync.WaitGroup, i *SysInfo) {
	i.User.Groups = runCmdAndGetOutput(1, "groups")
	defer wg.Done()
}
