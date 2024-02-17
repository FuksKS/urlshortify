package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

const (
	defaultAddr    = "localhost:8080"
	defaultBaseURL = "http://localhost:8000/qsd54gFg"
)

type Config struct {
	HTTPAddr string `env:"SERVER_ADDRESS"`
	BaseURL  string `env:"BASE_URL"`
}

func InitConfig() *Config {
	var cfg Config
	if err := envConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	flagAddr, flagBaseURL := flagConfig()

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = flagAddr
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = flagBaseURL
	}

	return &cfg
}

func flagConfig() (flagAddr, flagBaseURL string) {
	flag.StringVar(&flagAddr, "a", defaultAddr, "адрес запуска HTTP-сервера")
	flag.StringVar(&flagBaseURL, "b", defaultBaseURL, "базовый адрес результирующего сокращенного URL")
	flag.Parse()
	return
}

func envConfig(cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}
