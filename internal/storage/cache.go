package storage

import (
	"context"
	"github.com/FuksKS/urlshortify/internal/models"
	"sync"
)

type Cache struct {
	Cache      map[string]models.URLInfo
	mapRWMutex *sync.RWMutex
}

func NewCache(cache map[string]models.URLInfo) *Cache {
	return &Cache{Cache: cache, mapRWMutex: &sync.RWMutex{}}
}

func (c *Cache) Save(_ context.Context, _ map[string]models.URLInfo) error {
	// Для имплементации. За 1 раз кэш сохраняем только в файл
	return nil
}

func (c *Cache) SaveOneURL(_ context.Context, info models.URLInfo) error {
	if _, ok := c.Cache[info.ShortURL]; !ok {
		c.mapRWMutex.Lock()
		c.Cache[info.ShortURL] = info
		c.mapRWMutex.Unlock()
	}

	return nil
}

func (c *Cache) SaveURLs(_ context.Context, urls []models.URLInfo) error {
	for i := range urls {
		if _, ok := c.Cache[urls[i].ShortURL]; !ok {
			c.mapRWMutex.Lock()
			c.Cache[urls[i].ShortURL] = urls[i]
			c.mapRWMutex.Unlock()
		}
	}
	return nil
}

func (c *Cache) ReadAll(_ context.Context) (map[string]models.URLInfo, error) {
	// Для имплементации
	return c.Cache, nil
}

func (c *Cache) GetLongURL(_ context.Context, shortURL string) (models.URLInfo, error) {
	return c.Cache[shortURL], nil
}

func (c *Cache) GetUsersURLs(_ context.Context, _ string) ([]models.URLInfo, error) {
	// Для имплементации. Забираем урлы по юзеру всегда из бд
	return nil, nil
}

func (c *Cache) PingDB(_ context.Context) error {
	// Для имплементации
	return nil
}
