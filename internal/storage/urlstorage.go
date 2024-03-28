package storage

import (
	"sync"
)

type Storage struct {
	Cache      map[string]string
	mapRWMutex *sync.RWMutex
	saver      saver
	reader     reader
}

func New(filePath string) (*Storage, error) {
	saver, err := newFileWriter(filePath)
	if err != nil {
		return nil, err
	}
	reader, err := newFileReader(filePath)
	if err != nil {
		return nil, err
	}

	cashe, err := reader.read()
	if err != nil {
		return nil, err
	}

	st := &Storage{Cache: cashe, mapRWMutex: &sync.RWMutex{}, saver: saver, reader: reader}

	return st, nil
}

func (s *Storage) SaveShortURL(shortURL, longURL string) {
	if _, ok := s.Cache[shortURL]; !ok {
		s.mapRWMutex.Lock()
		s.Cache[shortURL] = longURL
		s.mapRWMutex.Unlock()
	}
}

func (s *Storage) GetLongURL(shortURL string) string {
	return s.Cache[shortURL]
}

func (s *Storage) SaveCache() error {
	if err := s.saver.Save(s.Cache); err != nil {
		return err
	}

	return nil
}
