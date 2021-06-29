package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// TODO: fix this mess

// FileUpload handles the /upload endpoint.
func FileUpload(w http.ResponseWriter, req *http.Request) {
	// retrieve the file from the request
	file, handler, err := req.FormFile("uploadFile")
	if err != nil {
		fmt.Printf("[-] Error retrieving file: %s\n", err)
		return
	}

	// read the file data
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("[-] Error reading the uploaded file: %s\n", err)
	}

	// create uploads dir if not existant yet
	_ = os.Mkdir("uploads", 0755)
	// read data into local file
	err = ioutil.WriteFile("uploads/"+handler.Filename, bytes, 0755)
	if err != nil {
		fmt.Printf("[-] Error creating and reading into local file: %s\n", err)
	}
	fmt.Println("[+] Successfully uploaded file")
	fmt.Fprintf(w, "Successfully uploaded file")
}

// FileDownload handles the /download endpoint.
func FileDownload(w http.ResponseWriter, req *http.Request) {
	// get the filename from the request
	filename := req.URL.Query().Get("file")
	if filename == "" {
		fmt.Println("[-] Download request doesn't contain file name")
		http.Error(w, "no file indicatd to download", 400)
		return
	}
	fmt.Println("[+] Payload wants to download ", filename)

	// open the file if it exists
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		fmt.Printf("[-] Error trying to open file: %s\n", err)
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
		fmt.Printf("[-] Error writing file into response: %s\n", err)
		return
	}
	fmt.Println("[+] Successfully downloaded file")
	fmt.Fprintf(w, "Successfully downloaded file")
}
