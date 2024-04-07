package storage

import (
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/pg"
	"github.com/google/uuid"
	"sync"
)

type Storage struct {
	Cache      map[string]string
	mapRWMutex *sync.RWMutex
	saver      saver
	reader     reader
}

func New(db pg.PgRepo, filePath string) (*Storage, error) {
	var saver saver
	var reader reader
	var err error

	if db.DB != nil {
		reader = &db
		saver = &db
	} else {
		saver, err = newFileWriter(filePath)
		if err != nil {
			return nil, fmt.Errorf("storage-New-newFileWriter-err: %w", err)
		}
		reader, err = newFileReader(filePath)
		if err != nil {
			return nil, fmt.Errorf("storage-New-newFileReader-err: %w", err)
		}
	}

	cashe, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("storage-New-reader-Read-err: %w", err)
	}

	st := &Storage{Cache: cashe, mapRWMutex: &sync.RWMutex{}, saver: saver, reader: reader}

	return st, nil
}

func (s *Storage) SaveShortURL(shortURL, longURL string) error {
	if _, ok := s.Cache[shortURL]; !ok {
		s.mapRWMutex.Lock()
		s.Cache[shortURL] = longURL
		s.mapRWMutex.Unlock()
	}

	err := s.saver.SaveOneURL(models.URLInfo{UUID: uuid.New().String(), ShortURL: shortURL, OriginalURL: longURL})

	return err
}

func (s *Storage) SaveURLs(urls []models.URLInfo) error {
	for i := range urls {
		if _, ok := s.Cache[urls[i].ShortURL]; !ok {
			s.mapRWMutex.Lock()
			s.Cache[urls[i].ShortURL] = urls[i].OriginalURL
			s.mapRWMutex.Unlock()
		}
	}

	return s.saver.SaveURLs(urls)
}

func (s *Storage) GetLongURL(shortURL string) string {
	return s.Cache[shortURL]
}

func (s *Storage) SaveCache() error {
	if err := s.saver.Save(s.Cache); err != nil {
		return fmt.Errorf("storage-SaveCache-Save-err: %w", err)
	}
	return nil
}
