package main

import (
	"fmt"
	"github.com/FuksKS/urlshortify/internal/config"
	"github.com/FuksKS/urlshortify/internal/handlers"
	"net/http"
)

func main() {

	cfg := config.InitConfig()

	fmt.Println(cfg.HTTPAddr)
	err := http.ListenAndServe(cfg.HTTPAddr, handlers.RootHandler(cfg.HTTPAddr, cfg.BaseURL))
	if err != nil {
		panic(err)
	}
}
