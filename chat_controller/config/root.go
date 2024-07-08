package config

import (
	"github.com/naoina/toml"
	"os"
)

type Config struct {
	DB struct {
		Database string
		Url      string
	}
	Kafka struct {
		Url     string
		GroupId string
	}

	Info struct {
		Port string
	}
}

func NewConfig(path string) *Config {
	config := new(Config)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	err = toml.NewDecoder(file).Decode(config)
	if err != nil {
		panic(err)
	} else {
		return config
	}
}
