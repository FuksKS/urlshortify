package handlers

import "github.com/FuksKS/urlshortify/internal/models"

type Storager interface {
	GetLongURL(shortURL string) string
	SaveShortURL(shortURL, longURL string)
	SaveURLs(urls []models.URLInfo)
	SaveDefaultURL(defaultURL, shortDefaultURL string)
	SaveCache() error
}
