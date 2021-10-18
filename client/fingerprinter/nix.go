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
	}

	// *nix only
	Distro string `json:"distro,omitempty"`
	Kernel string `json:"kernel,omitempty"`
}

func (i *SysInfo) fingerprint(wg *sync.WaitGroup) {
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
