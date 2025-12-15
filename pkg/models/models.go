package models

type Manga struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	Genres       []string `json:"genres"`
	Status       string   `json:"status"`
	TotalChapters int     `json:"total_chapters"`
	Description  string   `json:"description"`
	CoverURL     string   `json:"cover_url"`
}

type ProgressUpdate struct {
	UserID    string `json:"user_id"`
	MangaID   string `json:"manga_id"`
	Chapter   int    `json:"chapter"`
	Timestamp int64  `json:"timestamp"`
}
