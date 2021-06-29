package client

type config struct {
	Host string
	Port int
}

// TODO: Obfuscate client creds with XOR or something
func LoadClientConfig() config {
	var cfg config

	// Insert here your configuration for the client.
	cfg.Host = "127.0.0.1"
	cfg.Port = 1337

	return cfg
}
