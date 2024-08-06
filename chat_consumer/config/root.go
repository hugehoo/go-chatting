package config

import (
	"github.com/naoina/toml"
	"os"
)

type Config struct {
	Mongo struct {
		Database string
		Url      string
	}

	Kafka struct {
		Url      string
		GroupId  string
		ClientId string
	}
}

func NewConfig(path string) *Config {
	cfg := new(Config)
	if open, err := os.Open(path); err != nil {
		panic(err)
	} else if err := toml.NewDecoder(open).Decode(cfg); err != nil {
		panic(err)
	} else {
		return cfg
	}
}
