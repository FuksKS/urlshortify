package main

import (
	"github.com/FuksKS/urlshortify/internal/handlers"
	"net/http"
)

func main() {

	handler := handlers.New()

	err := http.ListenAndServe(handler.HTTPAddr, handler.RootHandler())
	if err != nil {
		panic(err)
	}
}
