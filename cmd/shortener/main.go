package main

import (
	"github.com/FuksKS/urlshortify/internal/handlers"
	"github.com/FuksKS/urlshortify/internal/storage"
	"net/http"
)

func main() {
	st := storage.New()

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.RootHandler(st))

	err := http.ListenAndServe(handlers.DefaultAddr, mux)
	if err != nil {
		panic(err)
	}
}
