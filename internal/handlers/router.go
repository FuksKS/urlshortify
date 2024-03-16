package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLHandler struct {
	storage  Storager
	db       *pgxpool.Pool
	HTTPAddr string
}

func New(st Storager, db *pgxpool.Pool, addr, baseURL string) *URLHandler {
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
	r.Get("/ping", withLogging(h.pingDb()))

	return r
}
