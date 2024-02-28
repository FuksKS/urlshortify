package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/logger"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/urlmaker"
	"go.uber.org/zap"
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

		logger.Log.Info("getURLID()", zap.String("incoming short URL", id), zap.String("longURL from storage", h.storage.GetLongURL(id)))

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

		logger.Log.Info("generateShortURL()", zap.String("incoming long URL", longURL), zap.String("short URL", shortURL))

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

		logger.Log.Info("shorten()", zap.String("incoming body", string(body)))
		req := models.ShortenReq{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			logger.Log.Info("shorten()", zap.String("Unmarshal body error", err.Error()))
			http.Error(w, "Unmarshal body error", http.StatusInternalServerError)
			return
		}

		shortURL := urlmaker.MakeShortURL(req.URL)
		h.storage.SaveShortURL(shortURL, req.URL)

		logger.Log.Info("shorten()", zap.String("incoming long URL", req.URL), zap.String("short URL", shortURL))

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
