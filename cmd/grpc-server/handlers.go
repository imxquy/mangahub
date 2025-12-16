package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	pb "mangahub/proto"
	tcpmod "mangahub/internal/tcp"
)

func (s *server) GetManga(ctx context.Context, req *pb.GetMangaRequest) (*pb.MangaResponse, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, title, author, genres, status, total_chapters, description
		FROM manga WHERE id = ?
	`, req.Id)

	var (
		id, title, author, genresText, status, desc string
		total int32
	)
	if err := row.Scan(&id, &title, &author, &genresText, &status, &total, &desc); err != nil {
		if err == sql.ErrNoRows {
			return &pb.MangaResponse{}, nil
		}
		return nil, err
	}

	genres := decodeGenres(genresText)
	return &pb.MangaResponse{
		Manga: &pb.Manga{
			Id:           id,
			Title:        title,
			Author:       author,
			Genres:       genres,
			Status:       status,
			TotalChapters: total,
			Description:  desc,
		},
	}, nil
}

func (s *server) SearchManga(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	q := strings.TrimSpace(req.Query)
	limit := req.Limit
	if limit <= 0 || limit > 50 { limit = 10 }

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, author, genres, status, total_chapters, description
		FROM manga
		WHERE title LIKE ? OR author LIKE ?
		LIMIT ?
	`, "%"+q+"%", "%"+q+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := &pb.SearchResponse{}
	for rows.Next() {
		var (
			id, title, author, genresText, status, desc string
			total int32
		)
		if err := rows.Scan(&id, &title, &author, &genresText, &status, &total, &desc); err != nil {
			return nil, err
		}
		out.Results = append(out.Results, &pb.Manga{
			Id:           id,
			Title:        title,
			Author:       author,
			Genres:       decodeGenres(genresText),
			Status:       status,
			TotalChapters: total,
			Description:  desc,
		})
	}
	return out, nil
}

func (s *server) UpdateProgress(ctx context.Context, req *pb.ProgressRequest) (*pb.ProgressResponse, error) {
	// update DB (upsert)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO user_progress (user_id, manga_id, current_chapter, status, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id, manga_id) DO UPDATE SET
			current_chapter=excluded.current_chapter,
			status=excluded.status,
			updated_at=CURRENT_TIMESTAMP
	`, req.UserId, req.MangaId, req.Chapter, "reading")
	if err != nil {
		return &pb.ProgressResponse{Ok: false}, err
	}

	// trigger TCP broadcast (same principle as HTTP)
	_ = s.tcpNotifier.SendProgressUpdate(tcpmod.ProgressUpdate{
		UserID:    req.UserId,
		MangaID:   req.MangaId,
		Chapter:   int(req.Chapter),
		Timestamp: time.Now().Unix(),
	})

	return &pb.ProgressResponse{Ok: true}, nil
}

func decodeGenres(genresText string) []string {
	genresText = strings.TrimSpace(genresText)
	if genresText == "" { return nil }

	// If stored as JSON array string (recommended)
	var arr []string
	if err := json.Unmarshal([]byte(genresText), &arr); err == nil {
		return arr
	}
	// fallback: comma-separated
	parts := strings.Split(genresText, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" { out = append(out, p) }
	}
	return out
}
