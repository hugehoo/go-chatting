package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	DB struct {
		Database string
		URL      string
	}
	Kafka struct {
		URL      string
		ClientID string
	}
}

func NewConfig(path string) *Config {
	config := new(Config)
	if f, err := os.Open(path); err != nil {
		panic(err)
	} else if err := toml.NewDecoder(f).Decode(config); err != nil {
		panic(err)
	} else {
		return config
	}
}
