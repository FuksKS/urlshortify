package storage

import (
	"sync"
)

type Storage struct {
	Cashe    map[string]string
	mapMutex *sync.Mutex
	saver    saver
	reader   reader
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

	st := &Storage{Cashe: cashe, mapMutex: &sync.Mutex{}, saver: saver, reader: reader}

	return st, nil
}

func (s *Storage) SaveShortURL(shortURL, longURL string) {
	if _, ok := s.Cashe[shortURL]; !ok {
		s.mapMutex.Lock()
		s.Cashe[shortURL] = longURL
		s.mapMutex.Unlock()
	}
}

func (s *Storage) GetLongURL(shortURL string) string {
	return s.Cashe[shortURL]
}

func (s *Storage) SaveDefaultURL(defaultURL, shortDefaultURL string) {
	s.mapMutex.Lock()
	s.Cashe[shortDefaultURL] = defaultURL
	s.mapMutex.Unlock()
}

func (s *Storage) SaveCache() error {
	if err := s.saver.Save(s.Cashe); err != nil {
		return err
	}

	return nil
}
