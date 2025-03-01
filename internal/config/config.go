package config

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port    string `env:"PORT" envDefault:"8080" yaml:"port"`
	Host    string `env:"HOST" envDefault:"localhost" yaml:"host"`
	Timeout int    `env:"TIMEOUT" envDefault:"5" yaml:"timeout"`
	Debug   bool   `env:"DEBUG" envDefault:"false" yaml:"debug"`
	Token   string `env:"TOKEN" envDefault:"" yaml:"token"`
	CloudID string `env:"CLOUD_ID" envDefault:"" yaml:"cloud_id"`
}

func New(configPath string) (*Config, error) {
	filename := filepath.Base(configPath)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	rawFileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config *Config
	err = yaml.Unmarshal(rawFileContent, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
