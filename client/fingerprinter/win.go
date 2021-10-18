//go:build windows
// +build windows

package fingerprinter

import (
	"sync"
)

type SysInfo struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	Name string `json:"name"`

	User struct {
		Name string `json:"name"`
	}
}

// holds all of the command functions to run
// TODO: make fingerprinter commands for windows
var commands []func(wg *sync.WaitGroup, i *SysInfo) = []func(wg *sync.WaitGroup, i *SysInfo){}
