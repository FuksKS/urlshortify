package pg

import (
	"context"
	"encoding/json"
	"github.com/FuksKS/urlshortify/internal/models"
	"time"
)

func (r *PgRepo) Save(cache map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	allURLs := make([]models.URLInfo, 0, len(cache))

	for shortURL, originalURL := range cache {
		allURLs = append(allURLs, models.URLInfo{
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

	return nil
}
