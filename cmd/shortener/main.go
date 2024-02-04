package main

import (
	"github.com/FuksKS/urlshortify/internal/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetURLID(w, r)
		case http.MethodPost:
			handlers.GenerateShortURL(w, r)
		default:
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(handlers.DefaultAddr, mux)
	if err != nil {
		panic(err)
	}
}
