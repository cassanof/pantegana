package server

import (
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/i582/cfmt/cmd/cfmt"
)

//go:generate go-bindata -o cert.go ../cert/...

type Listener struct {
	Cfg    *ListenerConfig
	Server *http.Server
}

type ListenerConfig struct {
	Addr      string
	Plaintext bool
	Verbose   bool
}

var listener *Listener = nil

func (cfg *ListenerConfig) SetupListener() *Listener {
	// read cert binary data from bundled assets
	certData, err := Asset("../cert/server.crt")
	if err != nil {
		cli.Print(cfmt.Sprintf("{{[-] Error reading cert file: %s\n}}::red", err))
	}
	// read key binary data from bundled assets
	keyData, err := Asset("../cert/server.key")
	if err != nil {
		cli.Print(cfmt.Sprintf("{{[-] Error reading cert file: %s\n}}::red", err))
	}

	// create the server with the custom pair
	cert, err := tls.X509KeyPair(certData, keyData)
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	server := http.Server{
		Addr:      cfg.Addr,
		TLSConfig: tlsConfig,
	}
	if cfg.Plaintext {
		server.TLSConfig = nil
	}

	return &Listener{
		Cfg:    cfg,
		Server: &server,
	}
}

func StartListener(cfg *ListenerConfig) {
	// check if a listener is already running
	if IsListening() {
		cli.Print(cfmt.Sprintln("{{[-] A listener is already running.}}::red"))
		return
	}

	// start the listener
	cli.Printf("[+] Listening on (%s)\n", cfg.Addr)
	listener = cfg.SetupListener()

	var err error
	if cfg.Plaintext {
		err = listener.Server.ListenAndServe()
	} else {
		err = listener.Server.ListenAndServeTLS("", "")
	}

	if err != http.ErrServerClosed {
		cli.PrintError(err)
		defer CloseListener()
	}
}

func CloseListener() error {
	var err error
	if listener != nil {
		err = listener.Server.Close()
		listener = nil
	} else {
		err = errors.New("There are no listeners running")
	}
	return err
}

func IsListening() bool {
	return listener != nil
}
