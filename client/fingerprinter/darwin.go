//go:build darwin
// +build darwin

package fingerprinter

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
		Name string `json:"name"`
	}
}

// TODO: fingerprintWindows
func (i *SysInfo) fingerprint(wg *sync.WaitGroup) {
}
