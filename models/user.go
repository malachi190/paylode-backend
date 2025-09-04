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
	Pin             string `json:"pin"`
	Password        string `json:"-"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func (m *UserModel) CreateUser(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO users 
	(first_name, last_name, email, email_verified_at, phone_number, password) 
	VALUES (?,?,?,?,?,?)
	`
	var result User

	res, err := m.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.EmailVerifiedAt, user.PhoneNumber, user.Password)

	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	err = m.DB.QueryRowContext(ctx, `SELECT id, first_name, last_name, email, phone_number, created_at, updated_at FROM users WHERE id = ?`, id).Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.PhoneNumber,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *UserModel) GetUserWithEmailOrPhone(email, phone string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, password, phone_number, created_at, updated_at
		FROM users 
		WHERE email = ?
		OR phone_number = ?
		LIMIT 1
	`

	var user User

	err := m.DB.QueryRowContext(ctx, query, email, phone).Scan(
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

func (m *UserModel) CreateUserPin(userID uint, pin string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET 
	pin = ?, 
	has_pin = ? 
	WHERE id = ?`

	err := m.DB.QueryRowContext(ctx, query, pin, true, userID).Err()

	if err != nil {
		return err
	}

	return nil
}
