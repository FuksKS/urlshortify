package pg

import (
	"context"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"time"
)

func (r *PgRepo) Save(_ context.Context, _ map[string]models.URLInfo) error {
	// Для имплементации. За 1 раз кэш сохраняем только в файл
	return nil
}

func (r *PgRepo) SaveOneURL(ctx context.Context, info models.URLInfo) error {
	ctx2, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	commandTag, err := r.DB.Exec(ctx2, saveOneURLQuery, info.UUID, info.ShortURL, info.OriginalURL, info.UserID)
	if err != nil {
		return fmt.Errorf("SaveOneURL-Exec-err: %w", err)
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrAffectNoRows
	}

	return nil
}

func (r *PgRepo) SaveURLs(ctx context.Context, urls []models.URLInfo) error {
	ctx2, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.DB.Begin(ctx2)
	if err != nil {
		tx.Rollback(ctx2)
		return fmt.Errorf("SaveURLs-BeginTx-err: %w", err)
	}
	defer tx.Rollback(ctx2)

	for i := range urls {
		_, err := tx.Exec(ctx2, saveOneURLQuery, urls[i].UUID, urls[i].ShortURL, urls[i].OriginalURL)
		if err != nil {
			tx.Rollback(ctx2)
			return fmt.Errorf("SaveURLs-saveOneURLQuery-Exec-err: %w", err)
		}
	}

	err = tx.Commit(ctx2)
	if err != nil {
		return fmt.Errorf("SaveURLs-Commit-err: %w", err)
	}

	return nil
}
