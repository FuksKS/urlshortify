package handlers

import "github.com/FuksKS/urlshortify/internal/models"

type Storager interface {
	GetLongURL(shortURL string) string
	SaveShortURL(info models.URLInfo) error
	SaveURLs(urls []models.URLInfo) error
	SaveCache() error
}
