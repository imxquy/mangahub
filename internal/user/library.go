package user

import (
	"time"
)

type LibraryItem struct {
	UserID        string    `json:"user_id"`
	MangaID       string    `json:"manga_id"`
	CurrentChapter int      `json:"current_chapter"`
	Status        string    `json:"status"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AddLibraryRequest struct {
	MangaID string `json:"manga_id"`
	Status  string `json:"status"` // reading/completed/plan_to_read...
}

func (r *Repo) AddToLibrary(userID, mangaID, status string) error {
	if status == "" {
		status = "reading"
	}
	_, err := r.DB.Exec(`
		INSERT INTO user_progress (user_id, manga_id, current_chapter, status, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id, manga_id) DO UPDATE SET
			status=excluded.status,
			updated_at=CURRENT_TIMESTAMP
	`, userID, mangaID, 1, status)
	return err
}

func (r *Repo) ListLibrary(userID string) ([]LibraryItem, error) {
	rows, err := r.DB.Query(`
		SELECT user_id, manga_id, current_chapter, status, updated_at
		FROM user_progress
		WHERE user_id = ?
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []LibraryItem
	for rows.Next() {
		var it LibraryItem
		if err := rows.Scan(&it.UserID, &it.MangaID, &it.CurrentChapter, &it.Status, &it.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, nil
}
