package handlers

import (
	"fmt"
	"github.com/FuksKS/urlshortify/internal/config"
	"github.com/FuksKS/urlshortify/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type URLHandler struct {
	storage  storage.Storager
	HttpAddr string
}

func New() *URLHandler {
	st := storage.New()

	cfg := config.InitConfig()

	st.SaveDefaultURL(cfg.HTTPAddr, cfg.BaseURL)

	return &URLHandler{
		storage:  st,
		HttpAddr: cfg.HTTPAddr,
	}
}

func (h *URLHandler) RootHandler() chi.Router {

	r := chi.NewRouter()

	r.Post("/", h.generateShortURL(h.HttpAddr))
	r.Get("/{id}", h.getURLID())

	return r
}

func (h *URLHandler) getURLID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/") // parts[0] == ""
		id := parts[1]

		if longURL := h.storage.GetLongURL(id); longURL != "" {
			http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
			return
		}

		// урла нет в хранилище
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func (h *URLHandler) generateShortURL(addr string) http.HandlerFunc {
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

		shortURL := h.storage.SaveShortURL(string(body))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "http://%s/%s", addr, shortURL)
	}
}
