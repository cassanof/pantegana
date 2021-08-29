// This file has some helper functions and utility functions related to the server.
package server

import (
	"net/http"
	"strings"
)

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")]
}
