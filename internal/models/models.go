package models

import "errors"

type ContextKey string

const UserIDKey ContextKey = "user_id"

var ErrAffectNoRows = errors.New("affect no rows")

type ShortenReq struct {
	URL string `json:"url"`
}

type ShortenResp struct {
	Result string `json:"result"`
}

type URLInfo struct {
	UUID        string `json:"correlation_id,omitempty" db:"id"`
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}
