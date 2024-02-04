package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	DefaultAddr = ":8080"
	defaultHost = "http://localhost" + DefaultAddr + "/"
)

var storage = make(map[string]string)

func GetURLID(w http.ResponseWriter, r *http.Request) {
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

	if longUrl, ok := storage[id]; ok {
		fmt.Println("longUrl: ", longUrl)
		http.Redirect(w, r, longUrl, http.StatusTemporaryRedirect)
		return
	}

	http.Error(w, "Bad request", http.StatusBadRequest)
}

func GenerateShortUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	/*
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" {
			http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
			return
		}
	*/

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading request body error", http.StatusInternalServerError)
		return
	}

	shortUrl := calculateHash(string(body))
	storage[shortUrl] = string(body)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s%s", defaultHost, shortUrl)
}
