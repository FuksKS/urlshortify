package handlers

import (
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	storage Storager
	BaseURL string
}

func New(st Storager, baseURL string) *URLHandler {
	return &URLHandler{
		storage: st,
		BaseURL: baseURL,
	}
}

func (h *URLHandler) InitRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(withLogging, withGzip)

	r.Post("/", h.shorten())
	r.Get("/{id}", h.getShorten())
	r.Post("/api/shorten", h.shortenJSON())

	return r
}
