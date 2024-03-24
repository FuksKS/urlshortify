package pg

import (
	"context"
	"fmt"
	"github.com/FuksKS/urlshortify/internal/models"
	"time"
)

func (r *PgRepo) Save(cache map[string]string) error {
	// Для имплементации. За 1 раз кэш сохраняем только в файл
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		allURLs := make([]models.URLInfo, 0, len(cache))

		for shortURL, originalURL := range cache {
			allURLs = append(allURLs, models.URLInfo{
				UUID:        uuid.New().String(),
				ShortURL:    shortURL,
				OriginalURL: originalURL,
			})
		}

		allURLsByte, err := json.Marshal(allURLs)
		if err != nil {
			return err
		}

		_, err = r.DB.Exec(ctx, saveCashQuery, string(allURLsByte))
		if err != nil {
			return err
		}

	*/

	return nil
}

func (r *PgRepo) SaveOneURL(info models.URLInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shortURL string
	err := r.DB.QueryRow(ctx, selectOneURLQuery).Scan(&shortURL)
	if err != nil {
		fmt.Println("SaveOneURL-selectOneURLQuery-err: ", err.Error())
		return err
	}

	fmt.Println("SaveOneURL-selectOneURLQuery-shortURL: ", shortURL)

	if shortURL != "" {
		return models.ErrAffectNoRows
	}

	_, err = r.DB.Exec(ctx, saveOneURLQuery, info.UUID, info.ShortURL, info.OriginalURL)
	if err != nil {
		fmt.Println("SaveOneURL-Exec-err: ", err.Error())
		return err
	}

	return nil
}

func (r *PgRepo) SaveURLs(urls []models.URLInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	defer tx.Rollback(ctx)

	for i := range urls {
		_, err := tx.Exec(ctx, saveOneURLQuery, urls[i].UUID, urls[i].ShortURL, urls[i].OriginalURL)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	return err
}
