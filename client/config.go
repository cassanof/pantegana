package client

type config struct {
	Host string
	Port int
}

func LoadClientConfig(host string, port int) *config {

	cfg := config{host, port}

	return &cfg
}
