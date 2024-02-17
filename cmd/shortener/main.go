package main

import (
	"github.com/FuksKS/urlshortify/internal/config"
	"github.com/FuksKS/urlshortify/internal/handlers"
	"github.com/FuksKS/urlshortify/internal/storage"
	"net/http"
)

func main() {

	st := storage.New()
	cfg := config.InitConfig()

	handler := handlers.New(st, cfg.HTTPAddr, cfg.HTTPAddr)

	err := http.ListenAndServe(handler.HTTPAddr, handler.InitRouter())
	if err != nil {
		panic(err)
	}
}
