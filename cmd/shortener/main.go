package main

import (
	"github.com/FuksKS/urlshortify/internal/config"
	"github.com/FuksKS/urlshortify/internal/handlers"
	"net/http"
)

func main() {

	cfg := config.InitConfig()

	err := http.ListenAndServe(cfg.HTTPAddr, handlers.RootHandler(cfg.HTTPAddr, cfg.BaseURL))
	if err != nil {
		panic(err)
	}
}
