package handlers

import (
	"io"
	"net/http"
	"strings"
)

const (
	DefaultAddr = ":8080"
	defaultHost = "http://localhost" + DefaultAddr + "/"
)

var storage = make(map[string]string)

func Shortify(GetUrlId http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			GetUrlId.ServeHTTP(w, r)
			return
		}

		/*
			contentType := r.Header.Get("Content-Type")
			if contentType != "text/plain" {
				http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
				return
			}

		*/

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		shortUrl := calculateHash(string(body))
		storage[shortUrl] = string(body)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte(defaultHost + shortUrl))
		return
	}
}

func GetUrlId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 2 { // parts[0] == ""
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id := parts[1]
	if id == "" { // Вызов был с эндпоинтом `/`. Ожидался метод POST
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusTemporaryRedirect)
	if longUrl, ok := storage[id]; ok {
		w.Header().Set("Location", longUrl)
		return
	}

	w.Header().Set("Location", defaultHost+id)

	return
}
