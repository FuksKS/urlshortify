package handlers

import (
	"errors"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/pg"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	storage  Storager
	db       pg.PgRepo
	HTTPAddr string
}

func New(st Storager, db pg.PgRepo, addr, baseURL string) (*URLHandler, error) {
	err := st.SaveShortURL(addr, baseURL)
	if err != nil && !errors.Is(err, models.ErrAffectNoRows) {
		return nil, fmt.Errorf("handlers-New-SaveDefaultURL-err: %w", err)
	}

	return &URLHandler{
		storage:  st,
		db:       db,
		HTTPAddr: addr,
	}, nil
}

func (h *URLHandler) InitRouter() chi.Router {

	r := chi.NewRouter()

	r.Post("/", withLogging(withGzip(h.generateShortURL())))
	r.Get("/{id}", withLogging(h.getURLID()))
	r.Post("/api/shorten", withLogging(withGzip(h.shorten())))
	r.Post("/api/shorten/batch", withLogging(withGzip(h.shortenBatch())))
	r.Get("/ping", withLogging(h.pingDB()))

	return r
}
