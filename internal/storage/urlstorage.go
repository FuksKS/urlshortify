package storage

import (
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
			return nil, err
		}
		reader, err = newFileReader(filePath)
		if err != nil {
			return nil, err
		}
	}

	cashe, err := reader.Read()
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

	s.saver.SaveURLs([]models.URLInfo{{uuid.New().String(), shortURL, longURL}})
}

func (s *Storage) SaveURLs(urls []models.URLInfo) {
	for i := range urls {
		if _, ok := s.Cashe[urls[i].ShortURL]; !ok {
			s.mapMutex.Lock()
			s.Cashe[urls[i].ShortURL] = urls[i].OriginalURL
			s.mapMutex.Unlock()
		}
	}

	s.saver.SaveURLs(urls)
}

func (s *Storage) GetLongURL(shortURL string) string {
	return s.Cashe[shortURL]
}

func (s *Storage) SaveDefaultURL(defaultURL, shortDefaultURL string) {
	s.mapMutex.Lock()
	s.Cashe[shortDefaultURL] = defaultURL
	s.mapMutex.Unlock()

	s.saver.SaveURLs([]models.URLInfo{{uuid.New().String(), shortDefaultURL, defaultURL}})
}

func (s *Storage) SaveCache() error {
	if err := s.saver.Save(s.Cashe); err != nil {
		return err
	}

	return nil
}
