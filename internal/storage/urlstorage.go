package storage

import (
	"crypto/sha256"
	"encoding/hex"
)

type Storage struct {
	cashe map[string]string
}

func New() *Storage {
	return &Storage{cashe: map[string]string{}}
}

func (s *Storage) SaveShortURL(input string) string {
	inputBytes := []byte(input)

	// Вычисление хэша с использованием SHA-256
	hash := sha256.Sum256(inputBytes)

	// Преобразование хэша в строку в шестнадцатеричном формате
	hashString := hex.EncodeToString(hash[:])

	shortURL := hashString[:8]
	if _, ok := s.cashe[shortURL]; !ok {
		s.cashe[shortURL] = input
	}

	return shortURL
}

func (s *Storage) GetLongURL(shortURL string) string {
	return s.cashe[shortURL]
}

type Storager interface {
	GetLongURL(shortURL string) string
	SaveShortURL(input string) string
}
