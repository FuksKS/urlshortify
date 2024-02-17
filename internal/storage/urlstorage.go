package storage

import (
	"sync"
)

type Storage struct {
	cashe    map[string]string
	mapMutex *sync.Mutex
}

func New() *Storage {
	return &Storage{cashe: map[string]string{}, mapMutex: &sync.Mutex{}}
}

func (s *Storage) SaveShortURL(shortURL, longURL string) {
	if _, ok := s.cashe[shortURL]; !ok {
		s.mapMutex.Lock()
		s.cashe[shortURL] = longURL
		s.mapMutex.Unlock()
	}

	return
}

func (s *Storage) GetLongURL(shortURL string) string {
	return s.cashe[shortURL]
}

func (s *Storage) SaveDefaultURL(defaultURL, shortDefaultURL string) {
	s.mapMutex.Lock()
	s.cashe[shortDefaultURL] = defaultURL
	s.mapMutex.Unlock()
}
