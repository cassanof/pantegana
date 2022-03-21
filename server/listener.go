package server

import (
	"crypto/tls"
	_ "embed"
	"errors"
	"io"
	"net/http"
)

type Listener struct {
	Cfg    *ListenerConfig
	Server *http.Server
}

type ListenerConfig struct {
	Addr      string
	Plaintext bool
	VW        io.Writer // set this to io.Discard if --verbose flag is off
}

var listener *Listener = nil

// Load the cert and key
//go:generate rm -fr cert
//go:generate mkdir cert
//go:generate cp ../cert/server.crt ./cert/
//go:generate cp ../cert/server.key ./cert/

//go:embed cert/server.crt
var certData []byte

//go:embed cert/server.key
var keyData []byte

func (cfg *ListenerConfig) SetupListener() *Listener {
	// create the server with the custom pair
	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		cli.Print(RedF("[-] Error pairing cert with key: %s\n", err))
	}
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
		cli.Print(RedF("[-] A listener is already running.\n"))
		return
	}

	// setup listener the listener
	listener = cfg.SetupListener()

	cli.Print(GreenF("[+] Listening on (%s)\n", cfg.Addr))

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
