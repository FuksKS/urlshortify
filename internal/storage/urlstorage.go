package storage

import (
	"context"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/FuksKS/urlshortify/internal/pg"
)

type Storage struct {
	cacheSaver  saver
	cacheReader reader
	saver       saver
	reader      reader
}

func New(ctx context.Context, filePath, DBDSN string) (*Storage, error) {
	var saver, cacheSaver saver
	var reader, cacheReader reader

	db, err := pg.NewConnect(ctx, DBDSN)
	if err != nil {
		fmt.Println("pg.NewConnect, DBDSN:", DBDSN, err.Error())
		return nil, fmt.Errorf("storage-New-pg.NewConnect-err: %w", err)
	}

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

	cache, err := reader.ReadAll(ctx)
	if err != nil {
		fmt.Println("reader.ReadAll:", err.Error())
		return nil, fmt.Errorf("storage-New-reader-Read-err: %w", err)
	}

	c := NewCache(cache)
	cacheSaver = c
	cacheReader = c

	st := &Storage{cacheSaver: cacheSaver, cacheReader: cacheReader, saver: saver, reader: reader}

	return st, nil
}

func (s *Storage) SaveShortURL(ctx context.Context, info models.URLInfo) error {
	err := s.cacheSaver.SaveOneURL(ctx, info)
	if err != nil {
		return fmt.Errorf("storage-SaveShortURL-cacheSaver-SaveOneURL-err: %w", err)
	}

	err = s.saver.SaveOneURL(ctx, info)
	if err != nil {
		return fmt.Errorf("storage-SaveShortURL-saver-SaveOneURL-err: %w", err)
	}

	return nil
}

func (s *Storage) SaveURLs(ctx context.Context, urls []models.URLInfo) error {
	err := s.cacheSaver.SaveURLs(ctx, urls)
	if err != nil {
		return fmt.Errorf("storage-SaveURLs-cacheSaver-SaveURLs-err: %w", err)
	}

	err = s.saver.SaveURLs(ctx, urls)
	if err != nil {
		return fmt.Errorf("storage-SaveURLs-saver-SaveURLs-err: %w", err)
	}

	return nil
}

func (s *Storage) GetLongURL(ctx context.Context, shortURL string) (models.URLInfo, error) {
	return s.cacheReader.GetLongURL(ctx, shortURL)
}

func (s *Storage) GetUsersURLs(ctx context.Context, userID string) ([]models.URLInfo, error) {
	return s.reader.GetUsersURLs(ctx, userID)
}

func (s *Storage) SaveCache(ctx context.Context) error {
	cache, err := s.cacheReader.ReadAll(ctx)
	if err != nil {
		return fmt.Errorf("storage-SaveCache-cacheReader-ReadAll-err: %w", err)
	}

	if err := s.saver.Save(ctx, cache); err != nil {
		return fmt.Errorf("storage-SaveCache-Save-err: %w", err)
	}
	return nil
}

func (s *Storage) PingDB(ctx context.Context) error {
	return s.reader.PingDB(ctx)
}

func (s *Storage) Shutdown(ctx context.Context) error {
	return s.reader.Shutdown(ctx)
}
