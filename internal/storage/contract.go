package storage

import (
	"context"
	"github.com/FuksKS/urlshortify/internal/models"
)

type saver interface {
	Save(ctx context.Context, cache map[string]models.URLInfo) error
	SaveOneURL(ctx context.Context, info models.URLInfo) error
	SaveURLs(ctx context.Context, urls []models.URLInfo) error
	Shutdown(ctx context.Context) error
}

type reader interface {
	ReadAll(ctx context.Context) (map[string]models.URLInfo, error)
	GetLongURL(ctx context.Context, shortURL string) (models.URLInfo, error)
	GetUsersURLs(ctx context.Context, userID string) ([]models.URLInfo, error)
	PingDB(ctx context.Context) error
}
