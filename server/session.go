package server

import (
	"errors"
	"net/http"
	"sync"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Session struct {
	Token   string
	Cmd     chan string
	Open    bool
	IP      string      // Could be both IPv4 and IPv6
	SysInfo interface{} // expects a struct, golang pls give me generics
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

	defer mutex.Unlock()
	return len(sessions) - 1, true
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

func ClearSessions() {
	sessions = []Session{}
}

func PrettyPrintSessions() {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"id", "ip", "hostname", "user", "os", "arch", "kernel", "token"})
	t.AppendSeparator()
	t.SetStyle(table.StyleColoredDark)
	for i, session := range sessions {
		if session.Open {
			sysInfoMap, _ := session.SysInfo.(map[string]interface{})
			userMap, _ := sysInfoMap["user"].(map[string]interface{})
			t.AppendRows([]table.Row{
				{
					i,
					session.IP,
					sysInfoMap["name"],
					userMap["name"],
					sysInfoMap["os"],
					sysInfoMap["arch"],
					sysInfoMap["kernel"],
					session.Token,
				},
			})
		}
	}
	cli.Printf("%s\n", t.Render())
}
