package handlers

type Storager interface {
	GetLongURL(shortURL string) string
	SaveShortURL(shortURL, longURL string)
	SaveDefaultURL(defaultURL, shortDefaultURL string)
}
