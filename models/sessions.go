package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type SessionModel struct {
	DB *sql.DB
}

type Sessions struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	TokenHash string `json:"token_hash"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

func (s *SessionModel) SaveRefreshToken(token string, userID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES (?,?,?)
	`
	ttl := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	_, err := s.DB.ExecContext(ctx, query, userID, token, ttl)

	if err != nil {
		return err
	}

	return nil
}

func (s *SessionModel) RefreshTokenExists(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT 1 FROM refresh_tokens WHERE token_hash = ? LIMIT 1`

	var dummy int

	err := s.DB.QueryRowContext(ctx, query, token).Scan(&dummy)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
