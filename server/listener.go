package server

import (
	"crypto/tls"
	"net/http"

	"github.com/i582/cfmt/cmd/cfmt"
)

//go:generate go-bindata -o cert.go ../cert/...

type listener struct {
	Cfg    *ListenerConfig
	Server *http.Server
}

type ListenerConfig struct {
	Addr      string
	Plaintext bool
	Verbose   bool
}

var Listener *listener

func SetupListener(cfg *ListenerConfig) *listener {
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

	return &listener{
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
	Listener = SetupListener(cfg)

	var err error
	if cfg.Plaintext {
		err = Listener.Server.ListenAndServe()
	} else {
		err = Listener.Server.ListenAndServeTLS("", "")
	}

	if err != http.ErrServerClosed {
		cli.PrintError(err)
		defer CloseListener()
	}
}
