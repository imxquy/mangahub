package user

import (
	"database/sql"
	"fmt"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

func (r *Repo) UpdateProgress(userID, mangaID string, chapter int, status string) error {
	if userID == "" || mangaID == "" {
		return fmt.Errorf("missing user_id or manga_id")
	}
	if chapter <= 0 {
		return fmt.Errorf("chapter must be > 0")
	}

	// Upsert into user_progress
	_, err := r.DB.Exec(`
		INSERT INTO user_progress (user_id, manga_id, current_chapter, status, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id, manga_id) DO UPDATE SET
			current_chapter=excluded.current_chapter,
			status=excluded.status,
			updated_at=CURRENT_TIMESTAMP
	`, userID, mangaID, chapter, status)
	return err
}
