package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB        *sql.DB
	JWTSecret string
}

func NewService(db *sql.DB, secret string) *Service {
	return &Service{DB: db, JWTSecret: secret}
}

func (s *Service) Register(username, email, password string) (userID string, err error) {
	if username == "" || password == "" {
		return "", errors.New("username/password required")
	}

	// Create deterministic-ish id (simple). Replace with uuid later if you want.
	userID = "usr_" + username

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// insert
	_, err = s.DB.Exec(`
		INSERT INTO users (id, username, email, password_hash)
		VALUES (?, ?, ?, ?)
	`, userID, username, nullIfEmpty(email), string(hash))
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *Service) Login(usernameOrEmail, password string) (token string, userID string, err error) {
	if usernameOrEmail == "" || password == "" {
		return "", "", errors.New("missing credentials")
	}

	var (
		id, username, email, passHash string
	)
	row := s.DB.QueryRow(`
		SELECT id, username, email, password_hash
		FROM users
		WHERE username = ? OR email = ?
	`, usernameOrEmail, usernameOrEmail)

	if err := row.Scan(&id, &username, &email, &passHash); err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub": id,
		"usr": username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	j := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := j.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", "", err
	}
	return tok, id, nil
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
