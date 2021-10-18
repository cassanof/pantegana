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

// TODO: fingerprintWindows
func (i *SysInfo) fingerprint(wg *sync.WaitGroup) {
}
