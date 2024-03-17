package handlers

import (
	"github.com/FuksKS/urlshortify/internal/pg"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	storage  Storager
	db       pg.PgRepo
	HTTPAddr string
}

func New(st Storager, db pg.PgRepo, addr, baseURL string) *URLHandler {
	st.SaveDefaultURL(addr, baseURL)

	return &URLHandler{
		storage:  st,
		db:       db,
		HTTPAddr: addr,
	}
}

func (h *URLHandler) InitRouter() chi.Router {

	r := chi.NewRouter()

	r.Post("/", withLogging(withGzip(h.generateShortURL())))
	r.Get("/{id}", withLogging(h.getURLID()))
	r.Post("/api/shorten", withLogging(withGzip(h.shorten())))
	r.Get("/ping", withLogging(h.pingDB()))

	return r
}
