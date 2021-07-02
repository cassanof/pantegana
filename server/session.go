package server

import (
	"errors"

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

func CreateSession(token string) int {
	// initialize sessions array
	if Sessions == nil {
		Sessions = make([]Session, 0)
	}

	index := FindSessionIndexByToken(token)
	if index != -1 {
		return index
	}

	session := Session{
		Token: token,
		Cmd:   make(chan string),
		Open:  true,
	}

	Sessions = append(Sessions, session)

	return len(Sessions) - 1
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
