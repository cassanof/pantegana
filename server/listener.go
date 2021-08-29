package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/i582/cfmt/cmd/cfmt"
)

//go:generate go-bindata -o cert.go ../cert/...

var Listener *http.Server

func SetupListener(host string, noTLS bool) *http.Server {
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
		Addr:      host,
		TLSConfig: tlsConfig,
	}
	if noTLS {
		server.TLSConfig = nil
	}

	return &server
}

func StartListener(host string, port int, noTLS bool) {
	// check if a listener is already running
	if IsListening() {
		cli.Print(cfmt.Sprintln("{{[-] A listener is already running.}}::red"))
		return
	}

	hoststr := fmt.Sprintf("%s:%d", host, port)

	// start the listener
	cli.Printf("[+] Listening on (%s)\n", hoststr)
	Listener = SetupListener(hoststr, noTLS)

	var err error
	if noTLS {
		err = Listener.ListenAndServe()
	} else {
		err = Listener.ListenAndServeTLS("", "")
	}

	if err != http.ErrServerClosed {
		cli.PrintError(err)
		CloseListener()
	}
}
