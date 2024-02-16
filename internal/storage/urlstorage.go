package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type Storager interface {
	GetLongURL(shortURL string) string
	SaveShortURL(input string) string
	SaveDefaultURL(defaultURL, shortDefaultURL string)
}

type Storage struct {
	cashe    map[string]string
	mapMutex *sync.Mutex
}

func New() *Storage {
	return &Storage{cashe: map[string]string{}, mapMutex: &sync.Mutex{}}
}

func (s *Storage) SaveShortURL(input string) string {
	inputBytes := []byte(input)

	// Вычисление хэша с использованием SHA-256
	hash := sha256.Sum256(inputBytes)

	// Преобразование хэша в строку в шестнадцатеричном формате
	hashString := hex.EncodeToString(hash[:])

	shortURL := hashString[:8]
	if _, ok := s.cashe[shortURL]; !ok {
		s.mapMutex.Lock()
		s.cashe[shortURL] = input
		s.mapMutex.Unlock()
	}

	return shortURL
}

func (s *Storage) GetLongURL(shortURL string) string {
	return s.cashe[shortURL]
}

func (s *Storage) SaveDefaultURL(defaultURL, shortDefaultURL string) {
	s.mapMutex.Lock()
	s.cashe[shortDefaultURL] = defaultURL
	s.mapMutex.Unlock()
}
