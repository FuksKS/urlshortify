package handlers

import (
	"fmt"
	"github.com/FuksKS/urlshortify/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

const (
	DefaultAddr = ":8080"
	defaultHost = "http://localhost" + DefaultAddr + "/"
)

/*
func RootHandler(storage storage.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getURLID(storage).ServeHTTP(w, r)
		case http.MethodPost:
			generateShortURL(storage).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

*/

func RootHandler() chi.Router {
	st := storage.New()

	r := chi.NewRouter()

	r.Post("/", generateShortURL(st))
	r.Get("/{id}", getURLID(st))

	return r
}

func getURLID(st storage.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if longURL := st.GetLongURL(id); longURL != "" {
			http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
			return
		}

		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func generateShortURL(st storage.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		shortURL := st.SaveShortURL(string(body))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s%s", defaultHost, shortURL)
	}
}
