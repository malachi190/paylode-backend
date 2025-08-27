package models

import (
	// "context"
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID              uint   `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	EmailVerifiedAt string `json:"email_verified_at"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"-"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func (m *UserModel) CreateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO users 
	(first_name, last_name, email, email_verified_at, phone_number, password) 
	VALUES (?,?,?,?,?,?)
	`
	_, err := m.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.EmailVerifiedAt, user.PhoneNumber, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) GetUserWithEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, password, phone_number, created_at, updated_at
		FROM users 
		WHERE email = ?
		LIMIT 1
	`

	var user User

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetUserById(id uint) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users 
		WHERE id = ?
		LIMIT 1
	`

	var user User

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
