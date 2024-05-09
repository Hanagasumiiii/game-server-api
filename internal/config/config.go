package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server   Server
	Postgres Postgres
}

type (
	Server struct {
		Adress string `yaml:"connection_string"`
	}
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}
)

func LoadConfig(filename string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
