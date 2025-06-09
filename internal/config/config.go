package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerURL    string        `yaml:"server_url"`
	PollInterval time.Duration `yaml:"poll_interval"`
	AuthToken    string        `yaml:"auth_token"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		PollInterval: 10 * time.Second,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
