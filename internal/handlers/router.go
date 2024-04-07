package handlers

import (
	"github.com/FuksKS/urlshortify/internal/pg"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	storage Storager
	db      pg.PgRepo
	BaseURL string
}

func New(st Storager, db pg.PgRepo, baseURL string) (*URLHandler, error) {
	return &URLHandler{
		storage: st,
		db:      db,
		BaseURL: baseURL,
	}, nil
}

func (h *URLHandler) InitRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(withLogging, withGzip)

	r.Post("/", h.shorten())
	r.Get("/{id}", h.getShorten())
	r.Post("/api/shorten", h.shortenJSON())
	r.Post("/api/shorten/batch", h.shortenBatch())
	r.Get("/ping", h.pingDB())

	return r
}
