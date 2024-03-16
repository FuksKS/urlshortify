package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
)

const (
	defaultAddr     = "localhost:8080"
	defaultBaseURL  = "http://localhost:8000/qsd54gFg"
	defaultFilePath = "/tmp/short-url-db.json"
)

const (
	defaultHost   = "localhost"
	defaultUser   = "defUser"
	defaultPass   = "pass"
	defaultDbName = "shortener"
)

type Config struct {
	HTTPAddr    string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
	DbDSN       string `env:"DATABASE_DSN"`
}

func Init() *Config {
	var cfg Config
	if err := envConfig(&cfg); err != nil {
		log.Fatal(err)
	}

	flagAddr, flagBaseURL, flagFilePath, flagDbDSN := flagConfig()

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = flagAddr
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = flagBaseURL
	}
	if cfg.FileStorage == "" {
		cfg.FileStorage = flagFilePath
	}
	if cfg.DbDSN == "" {
		cfg.DbDSN = flagDbDSN
	}

	return &cfg
}

func flagConfig() (flagAddr, flagBaseURL, flagFilePath, flagDbDSN string) {
	flag.StringVar(&flagAddr, "a", defaultAddr, "адрес запуска HTTP-сервера")
	flag.StringVar(&flagBaseURL, "b", defaultBaseURL, "базовый адрес результирующего сокращенного URL")
	flag.StringVar(&flagFilePath, "f", defaultFilePath, "полное имя файла, куда сохраняются данные в формате JSON")

	defaultDbDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", defaultHost, defaultUser, defaultPass, defaultDbName)
	flag.StringVar(&flagDbDSN, "d", defaultDbDSN, "строка с адресом подключения к БД")
	flag.Parse()
	return
}

func envConfig(cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}
