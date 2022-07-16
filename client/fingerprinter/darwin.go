//go:build darwin
// +build darwin

package fingerprinter

type SysInfo struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	Name string `json:"name"`

	User struct {
		Name string `json:"name"`
	} `json:"user"`
}

// holds all of the command functions to run
// TODO: make fingerprinter commands for osx
var commands []command = []command{}
