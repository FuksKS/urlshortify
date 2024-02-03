package main

import (
	"github.com/FuksKS/urlshortify/internal/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, handlers.Shortify(http.HandlerFunc(handlers.GetUrlId)))

	err := http.ListenAndServe(handlers.DefaultAddr, mux)
	if err != nil {
		panic(err)
	}
}
