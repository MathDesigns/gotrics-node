package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerAddress string `yaml:"server_address"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf(err.Error())
		}
	}(file)

	var cfg Config
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
