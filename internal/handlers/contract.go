package handlers

import "github.com/FuksKS/urlshortify/internal/models"

type Storager interface {
	GetLongURL(shortURL string) string
	SaveShortURL(shortURL, longURL string)
	SaveDefaultURL(defaultURL, shortDefaultURL string)
}

type FileWriter interface {
	WriteToFile(info models.URLInfo) error
}

type FileReader interface {
	ReadFromFile(shortURL string) (string, error)
}
