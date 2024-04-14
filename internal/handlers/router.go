package handlers

import (
	"github.com/FuksKS/urlshortify/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	storage Storager
	BaseURL string
}

func New(st Storager, baseURL string) (*URLHandler, error) {
	return &URLHandler{
		storage: st,
		BaseURL: baseURL,
	}, nil
}

func (h *URLHandler) InitRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.WithLogging, middleware.WithAuth, middleware.WithGzip)

	r.Post("/", h.shorten())
	r.Get("/{id}", h.getShorten())
	r.Post("/api/shorten", h.shortenJSON())
	r.Post("/api/shorten/batch", h.shortenBatch())
	r.Get("/api/user/urls", h.getUsersShorten())
	r.Get("/ping", h.pingDB())

	return r
}
