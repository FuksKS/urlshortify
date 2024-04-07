package pg

import (
	"context"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/jackc/pgx/v5"
	"time"
)

func (r *PgRepo) ReadAll(ctx context.Context) (map[string]models.URLInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx, getAllURLsQuery)
	if err != nil {
		return nil, fmt.Errorf("PgRepo-Read-Query-err: %w", err)
	}
	defer rows.Close()

	urlsInfo, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.URLInfo])
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("PgRepo-Read-CollectRows-err: %w", err)
	}

	m := make(map[string]models.URLInfo, len(urlsInfo))
	for _, urlInfo := range urlsInfo {
		m[urlInfo.ShortURL] = urlInfo
	}

	return m, nil
}

func (r *PgRepo) GetUsersURLs(ctx context.Context, userID string) ([]models.URLInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx, getUsersURLsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("PgRepo-GetUsersURLs-Query-err: %w", err)
	}
	defer rows.Close()

	urlsInfo, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.URLInfo])
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("PgRepo-GetUsersURLs-CollectRows-err: %w", err)
	}

	return urlsInfo, nil
}

func (r *PgRepo) GetLongURL(_ context.Context, _ string) (models.URLInfo, error) {
	// Для имплементации. Один урл берем только из кэша
	return models.URLInfo{}, nil
}

func (r *PgRepo) PingDB(_ context.Context) error {
	// Для имплементации
	return nil
}
