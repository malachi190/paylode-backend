package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type CardModel struct {
	DB *sql.DB
}

type Card struct {
	ID             uint   `json:"id"`
	UserID         uint   `json:"user_id"`
	Brand          string `json:"brand"`
	LastFourDigits string `json:"last4"`
	Token          string `json:"token"`
	ExpiryMonth    string `json:"exp_month"`
	ExpiryYear     string `json:"exp_year"`
	CreatedAt      string `json:"created_at"`
	User           *User  `json:"user,omitempty"`
}

func (c *CardModel) CreateCard(card *Card) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO cards (user_id, brand, last4, token, exp_month, exp_year) VALUES (?,?,?,?,?,?)`

	_, err := c.DB.ExecContext(ctx, query, card.UserID, card.Brand, card.LastFourDigits, card.Token, card.ExpiryMonth, card.ExpiryYear)

	if err != nil {
		return err
	}

	return nil
}

func (c *CardModel) GetCards(userID uint) ([]Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, brand, last4, token, exp_month, exp_year, created_at FROM cards WHERE user_id = ?`

	var cards []Card

	rows, err := c.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Card
		err := rows.Scan(&c.ID, &c.UserID, &c.Brand, &c.LastFourDigits, &c.Token, &c.ExpiryMonth, &c.ExpiryYear, &c.CreatedAt)

		if err != nil {
			return nil, err
		}

		cards = append(cards, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (c *CardModel) ValidateCardToken(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT 1 FROM cards WHERE token = ? LIMIT 1`

	var dummy int

	err := c.DB.QueryRowContext(ctx, query, token).Scan(&dummy)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
