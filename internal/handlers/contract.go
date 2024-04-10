package handlers

import "github.com/FuksKS/urlshortify/internal/models"

type Storager interface {
	GetLongURL(shortURL string) string
	GetUsersURLs(userID string) ([]models.URLInfo, error)
	SaveShortURL(models.URLInfo) error
	SaveURLs(urls []models.URLInfo) error
	SaveCache() error
}
