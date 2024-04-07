package pg

import (
	"context"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"github.com/jackc/pgx/v5"
	"time"
)

func (r *PgRepo) Read() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

	m := make(map[string]string, len(urlsInfo))
	for _, urls := range urlsInfo {
		m[urls.ShortURL] = urls.OriginalURL
	}

	return m, nil
}
