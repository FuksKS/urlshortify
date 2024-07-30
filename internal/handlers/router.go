package handlers

import (
	"github.com/FuksKS/urlshortify/internal/middleware"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	storage Storager
	BaseURL string

	// канал для сбора урлов на удаление
	DeleteURLChan chan models.DeleteURLs
}

func New(st Storager, baseURL string) (*URLHandler, error) {
	urlHandler := &URLHandler{
		storage:       st,
		BaseURL:       baseURL,
		DeleteURLChan: make(chan models.DeleteURLs, 1024), // установим каналу буфер в 1024 сообщения
	}

	// запустим горутину с фоновым сохранением урлов для удаления
	go urlHandler.flushDeleteURLs()

	return urlHandler, nil
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
	r.Delete("/api/user/urls", h.deleteShortenBatch())

	return r
}
