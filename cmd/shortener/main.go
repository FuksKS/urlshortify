package main

import (
	"github.com/FuksKS/urlshortify/internal/handlers"
	"net/http"
)

func main() {

	err := http.ListenAndServe(handlers.DefaultAddr, handlers.RootHandler())
	if err != nil {
		panic(err)
	}
}
