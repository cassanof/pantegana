package main

import (
	client "github.com/elleven11/pantegana/client"
)

func main() {

	cfg := client.ClientConfig{
		Host:    "127.0.0.1",
		Port:    1337,
		HasTLS:  true, // for debug only
		HasLogs: true, // disable this in "production"
	}

	client.RunClient(&cfg)
}
