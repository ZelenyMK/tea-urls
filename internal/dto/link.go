package dto

type Request struct {
	URL string `form:"url"`
}

type Response struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortenURL  string `json:"short_url"`
	Alias       string `json:"alias"`
}
