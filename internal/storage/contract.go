package storage

import (
	"github.com/FuksKS/urlshortify/internal/models"
)

type saver interface {
	Save(cache map[string]string) error
	SaveOneURL(info models.URLInfo) error
	SaveURLs(urls []models.URLInfo) error
}

type reader interface {
	Read() (map[string]string, error)
}
