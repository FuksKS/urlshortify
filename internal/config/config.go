package config

import "flag"

type Config struct {
	HTTPAddr string
	BaseURL  string
}

func InitConfig() *Config {
	httpAddr := flag.String("a", "localhost:8888", "адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "http://localhost:8000/qsd54gFg", "базовый адрес результирующего сокращенного URL")
	flag.Parse()

	config := &Config{
		HTTPAddr: *httpAddr,
		BaseURL:  *baseURL,
	}

	return config
}
