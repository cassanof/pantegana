package server

type Session struct {
	Token string
	Cmd   chan string
	Open  bool
}

var Sessions []Session

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

func FindSessionIndexByToken(token string) int {
	for i := 0; i < len(Sessions); i++ {
		if Sessions[i].Token == token {
			return i
		}
	}
	return -1
}
