package manga

import (
	"database/sql"
	"encoding/json"
	"strings"
)

type Repo struct{ DB *sql.DB }

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

type MangaRow struct {
	ID            string
	Title         string
	Author        string
	GenresText    string
	Status        string
	TotalChapters int
	Description   string
	CoverURL      string
}

func (r *Repo) GetByID(id string) (*Manga, error) {
	row := r.DB.QueryRow(`
		SELECT id, title, author, genres, status, total_chapters, description, cover_url
		FROM manga WHERE id = ?
	`, id)

	var m MangaRow
	if err := row.Scan(&m.ID, &m.Title, &m.Author, &m.GenresText, &m.Status, &m.TotalChapters, &m.Description, &m.CoverURL); err != nil {
		return nil, err
	}
	return mapRow(m), nil
}

func (r *Repo) Search(query, genre, status string, limit int) ([]Manga, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	q := "%" + strings.TrimSpace(query) + "%"
	st := strings.TrimSpace(status)

	// Simplest: filter by title/author + status. Genre filtering will be string match on JSON text.
	genreLike := "%" + strings.TrimSpace(genre) + "%"

	rows, err := r.DB.Query(`
		SELECT id, title, author, genres, status, total_chapters, description, cover_url
		FROM manga
		WHERE (title LIKE ? OR author LIKE ?)
		  AND (? = '' OR status = ?)
		  AND (? = '' OR genres LIKE ?)
		LIMIT ?
	`, q, q, st, st, strings.TrimSpace(genre), genreLike, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Manga
	for rows.Next() {
		var m MangaRow
		if err := rows.Scan(&m.ID, &m.Title, &m.Author, &m.GenresText, &m.Status, &m.TotalChapters, &m.Description, &m.CoverURL); err != nil {
			return nil, err
		}
		out = append(out, *mapRow(m))
	}
	return out, nil
}

func mapRow(rw MangaRow) *Manga {
	var genres []string
	_ = json.Unmarshal([]byte(rw.GenresText), &genres)
	return &Manga{
		ID:            rw.ID,
		Title:         rw.Title,
		Author:        rw.Author,
		Genres:        genres,
		Status:        rw.Status,
		TotalChapters: rw.TotalChapters,
		Description:   rw.Description,
		CoverURL:      rw.CoverURL,
	}
}
