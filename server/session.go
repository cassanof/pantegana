package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	client "github.com/elleven11/pantegana/client"
)

type Session struct {
	Token   string
	Cmd     chan string
	Open    bool
	SysInfo client.SysInfo
}

var Sessions []Session

// errors
var ErrSessionIsNotOpen = errors.New("The requested session is not open.")

func CreateSession(token string) (int, bool) {
	// initialize sessions array
	if Sessions == nil {
		Sessions = make([]Session, 0)
	}

	index := FindSessionIndexByToken(token)
	if index != -1 {
		return index, false
	}

	session := Session{
		Token: token,
		Cmd:   make(chan string),
		Open:  true,
	}

	Sessions = append(Sessions, session)

	return len(Sessions) - 1, true
}

func (s *Session) WriteToCmd(command string) error {
	if s.Open == false {
		return ErrSessionIsNotOpen
	}
	s.Cmd <- command
	return nil
}

func FindSessionIndexByToken(token string) int {
	for i := 0; i < len(Sessions); i++ {
		if Sessions[i].Token == token {
			return i
		}
	}
	return -1
}

func PrettyPrintSessions() {
	header := "||                          Sessions                          ||"
	spacer := strings.Repeat("=", len(header))
	output := fmt.Sprintf("%s\n%s\n%s\n", spacer, header, spacer)
	for i, session := range Sessions {
		if session.Open {
			fragment := fmt.Sprintf("|| ID: %d - Token: %s", i, session.Token)
			sessionInfo := fmt.Sprintf("%s%s||\n%s\n", fragment, strings.Repeat(" ", len(header)-len(fragment)-2), spacer)
			json, _ := json.MarshalIndent(session.SysInfo, "||", "\t")
			sessionInfo = fmt.Sprintf("%s||%s\n%s\n", sessionInfo, json, spacer)
			output += sessionInfo
		}
	}
	cli.Printf("%s", output)
}
