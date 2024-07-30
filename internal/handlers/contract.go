package handlers

import (
	"context"
	"github.com/FuksKS/urlshortify/internal/models"
)

type Storager interface {
	GetLongURL(ctx context.Context, shortURL string) (models.URLInfo, error)
	GetUsersURLs(ctx context.Context, userID string) ([]models.URLInfo, error)
	SaveShortURL(ctx context.Context, info models.URLInfo) error
	SaveURLs(ctx context.Context, urls []models.URLInfo) error
	SaveCache(ctx context.Context) error
	DeleteURLs(ctx context.Context, deleteURLs []models.DeleteURLs) error
	PingDB(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
