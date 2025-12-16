package manga

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

type Manga struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Genres        []string `json:"genres"`
	Status        string   `json:"status"`
	TotalChapters int      `json:"total_chapters"`
	Description   string   `json:"description"`
	CoverURL      string   `json:"cover_url"`
}

func SeedIfEmpty(db *sql.DB, jsonPath string) error {
	// 1) check count
	var count int
	if err := db.QueryRow(`SELECT COUNT(1) FROM manga`).Scan(&count); err != nil {
		return fmt.Errorf("count manga: %w", err)
	}
	if count > 0 {
		return nil
	}

	// 2) read file
	b, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("read manga json: %w", err)
	}

	var list []Manga
	if err := json.Unmarshal(b, &list); err != nil {
		return fmt.Errorf("unmarshal manga json: %w", err)
	}
	if len(list) == 0 {
		return fmt.Errorf("manga.json empty")
	}

	// 3) insert in tx
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO manga (id, title, author, genres, status, total_chapters, description, cover_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, m := range list {
		if m.ID == "" || m.Title == "" {
			continue
		}
		genresText, _ := json.Marshal(m.Genres) // store as JSON string in TEXT column
		if _, err := stmt.Exec(
			m.ID,
			m.Title,
			m.Author,
			string(genresText),
			m.Status,
			m.TotalChapters,
			m.Description,
			m.CoverURL,
		); err != nil {
			return fmt.Errorf("insert manga %s: %w", m.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
