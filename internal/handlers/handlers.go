package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
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

func (h *URLHandler) shorten() http.HandlerFunc {
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

		fullHost := fmt.Sprintf("http://%s/%s", h.HTTPAddr, shortURL)
		resp := models.ShortenResp{Result: fullHost}
		respB, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(respB)
	}
}

func (h *URLHandler) shortenBatch() http.HandlerFunc {
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

		var req []models.URLInfo
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Unmarshal body error", http.StatusInternalServerError)
			return
		}

		for i := range req {
			req[i].ShortURL = urlmaker.MakeShortURL(req[i].OriginalURL)
		}

		h.storage.SaveURLs(req)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		resp := make([]models.URLInfo, len(req))
		for i := range req {
			resp[i].UUID = req[i].UUID
			resp[i].ShortURL = fmt.Sprintf("http://%s/%s", h.HTTPAddr, req[i].ShortURL)
		}

		respB, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Marshal response error", http.StatusInternalServerError)
			return
		}

		w.Write(respB)
	}
}

func (h *URLHandler) pingDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := h.db.DB.Ping(context.Background()); err != nil {
			http.Error(w, "Ping db", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
