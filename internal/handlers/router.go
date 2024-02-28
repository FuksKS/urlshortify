package handlers

import (
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	storage  Storager
	HTTPAddr string
}

func New(st Storager, addr, baseURL string) *URLHandler {
	st.SaveDefaultURL(addr, baseURL)

	return &URLHandler{
		storage:  st,
		HTTPAddr: addr,
	}
}

func (h *URLHandler) InitRouter() chi.Router {

	r := chi.NewRouter()

	r.Post("/", withLogging(withGzip(h.generateShortURL())))
	r.Get("/{id}", withLogging(h.getURLID()))
	r.Post("/api/shorten", withLogging(withGzip(h.shorten())))

	return r
}
