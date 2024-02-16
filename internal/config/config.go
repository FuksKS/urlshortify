package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

const (
	DefaultAddr    = "localhost:8080"
	defaultBaseURL = "http://localhost:8000/qsd54gFg"
)

type Config struct {
	HTTPAddr string `env:"SERVER_ADDRESS"`
	BaseURL  string `env:"BASE_URL"`
}

func InitConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	var flagAddr, flagBaseURL string
	flag.StringVar(&flagAddr, "a", DefaultAddr, "адрес запуска HTTP-сервера")
	flag.StringVar(&flagBaseURL, "b", defaultBaseURL, "базовый адрес результирующего сокращенного URL")
	flag.Parse()

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = flagAddr
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = flagBaseURL
	}

	return &cfg
}
