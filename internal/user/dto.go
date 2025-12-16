package user

type UpdateProgressRequest struct {
	MangaID  string `json:"manga_id"`
	Chapter  int    `json:"chapter"`
	Status   string `json:"status"`
}
