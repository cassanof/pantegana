package main

import (
	client "github.com/elleven11/pantegana/client"
)

func main() {
	config := client.LoadClientConfig("127.0.0.1", 1337)

	client.RunClient(config.Host, config.Port)
}
