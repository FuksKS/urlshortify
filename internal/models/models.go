package models

type ShortenReq struct {
	URL string `json:"url"`
}

type ShortenResp struct {
	Result string `json:"result"`
}
