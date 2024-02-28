package models

type ShortenReq struct {
	URL string `json:"url"`
}

type ShortenResp struct {
	Result string `json:"result"`
}

type URLInfo struct {
	Uuid        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
