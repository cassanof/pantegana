// This file has some helper functions and utility functions related to the server.
package server

import (
	"errors"
	"net/http"
	"strings"
)

func CloseListener() error {
	var err error
	if Listener != nil {
		err = Listener.Close()
		Listener = nil
	} else {
		err = errors.New("There are not listeners running")
	}
	return err
}

func IsListening() bool {
	if Listener != nil {
		return true
	} else {
		return false
	}
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return ip
}
