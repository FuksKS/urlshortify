package handlers

import (
	"fmt"
	"github.com/FuksKS/urlshortify/internal/urlmaker"
	"io"
	"net/http"
	"strings"
)

func (h *URLHandler) getURLID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/") // parts[0] == ""
		if len(parts) != 2 {
			http.Error(w, "Incorrect path", http.StatusBadRequest)
		}
		id := parts[1]

		if longURL := h.storage.GetLongURL(id); longURL != "" {
			http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
			return
		}

		// урла нет в хранилище
		http.Error(w, "Unknown short URL", http.StatusBadRequest)
	}
}

func (h *URLHandler) generateShortURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		longURL := string(body)
		shortURL := urlmaker.MakeShortURL(longURL)
		h.storage.SaveShortURL(shortURL, longURL)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		resp := fmt.Sprintf("http://%s/%s", h.HTTPAddr, shortURL)
		w.Write([]byte(resp))
	}
}
