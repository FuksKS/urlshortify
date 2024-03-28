package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/urlmaker"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

func (h *URLHandler) getShorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			// при локальных тестах chi.RouteContext(r.Context()) = nil
			parts := strings.Split(r.URL.Path, "/") // parts[0] == ""
			if len(parts) != 2 {
				http.Error(w, "Incorrect path", http.StatusBadRequest)
			}
			id = parts[1]
		}

		if longURL := h.storage.GetLongURL(id); longURL != "" {
			http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
			return
		}

		// урла нет в хранилище
		http.Error(w, "Unknown short URL", http.StatusBadRequest)
	}
}

func (h *URLHandler) shorten() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		resp := fmt.Sprintf("%s%s", h.BaseURL, shortURL)
		w.Write([]byte(resp))
	}
}

func (h *URLHandler) shortenJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading request body error", http.StatusInternalServerError)
			return
		}

		req := models.ShortenReq{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Unmarshal body error", http.StatusInternalServerError)
			return
		}

		shortURL := urlmaker.MakeShortURL(req.URL)
		h.storage.SaveShortURL(shortURL, req.URL)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fullHost := fmt.Sprintf("%s%s", h.BaseURL, shortURL)
		resp := models.ShortenResp{Result: fullHost}
		respB, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(respB)
	}
}
