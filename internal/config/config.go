package config

import (
	"flag"
	"os"
)

const (
	defaultPort    = 8080
	defaultHost    = "0.0.0.0"
	defaultTimeout = 5
	defaultDebug   = false
)

type Config struct {
	Port    uint16
	Host    string
	Timeout int
	Debug   bool
	Token   string
	CloudID string
}

func New() (*Config, error) {
	config := initConfig()
	return config, nil
}

func initConfig() *Config {
	port := flag.Uint("port", defaultPort, "Port to listen on")
	host := flag.String("host", defaultHost, "Host to listen on")
	timeout := flag.Int("timeout", defaultTimeout, "Request timeout in seconds")
	debug := flag.Bool("debug", defaultDebug, "Enable debug mode")

	flag.Parse()

	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is required")
	}

	cloudID := os.Getenv("CLOUD_ID")
	if cloudID == "" {
		panic("CLOUD_ID environment variable is required")
	}

	return &Config{
		Port:    uint16(*port),
		Host:    *host,
		Timeout: *timeout,
		Debug:   *debug,
		Token:   token,
		CloudID: cloudID,
	}
}
