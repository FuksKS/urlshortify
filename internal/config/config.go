package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

const (
	defaultAddr     = "localhost:8080"
	defaultBaseURL  = "http://localhost:8080/"
	defaultFilePath = "/tmp/short-url-db.json"
)

type Config struct {
	HTTPAddr    string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
}

func Init() *Config {
	var cfg Config
	if err := envConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	flagAddr, flagBaseURL, flagFilePath := flagConfig()

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = flagAddr
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = flagBaseURL
	}
	if cfg.FileStorage == "" {
		cfg.FileStorage = flagFilePath
	}

	return &cfg
}

func flagConfig() (flagAddr, flagBaseURL, flagFilePath string) {
	flag.StringVar(&flagAddr, "a", defaultAddr, "адрес запуска HTTP-сервера")
	flag.StringVar(&flagBaseURL, "b", defaultBaseURL, "базовый адрес результирующего сокращенного URL")
	flag.StringVar(&flagFilePath, "f", defaultFilePath, "полное имя файла, куда сохраняются данные в формате JSON")
	flag.Parse()
	return
}

func envConfig(cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}
