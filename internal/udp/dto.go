package udp

type NotifyRequest struct {
	MangaID string `json:"manga_id"`
	Message string `json:"message"`
}
