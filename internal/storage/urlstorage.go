package storage

import (
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/pg"
	"github.com/google/uuid"
	"sync"
)

type Storage struct {
	Cashe    map[string]string
	mapMutex *sync.Mutex
	saver    saver
	reader   reader
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

	st := &Storage{Cashe: cashe, mapMutex: &sync.Mutex{}, saver: saver, reader: reader}

	return st, nil
}

func (s *Storage) SaveShortURL(shortURL, longURL string) error {
	if _, ok := s.Cashe[shortURL]; !ok {
		s.mapMutex.Lock()
		s.Cashe[shortURL] = longURL
		s.mapMutex.Unlock()
	}

	allURLs := make([]string, 0, len(s.Cashe))
	for shURL := range s.Cashe {
		allURLs = append(allURLs, shURL)
	}

	err := s.saver.SaveOneURL(models.URLInfo{UUID: uuid.New().String(), ShortURL: shortURL, OriginalURL: longURL})

	return err
}

func (s *Storage) SaveURLs(urls []models.URLInfo) error {
	for i := range urls {
		if _, ok := s.Cashe[urls[i].ShortURL]; !ok {
			s.mapMutex.Lock()
			s.Cashe[urls[i].ShortURL] = urls[i].OriginalURL
			s.mapMutex.Unlock()
		}
	}

	return s.saver.SaveURLs(urls)
}

func (s *Storage) GetLongURL(shortURL string) string {
	return s.Cashe[shortURL]
}

func (s *Storage) SaveCache() error {
	if err := s.saver.Save(s.Cashe); err != nil {
		return fmt.Errorf("storage-SaveCache-Save-err: %w", err)
	}

	return nil
}
