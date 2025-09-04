package models

import (
	"context"
	"database/sql"
	"time"
)

type TransactionModel struct {
	DB *sql.DB
}

type Transaction struct {
	ID                   uint    `json:"id"`
	UserID               uint    `json:"user_id"`
	TransactionType      string  `json:"transaction_type"`
	TransactionReference string  `json:"reference"`
	Amount               float64 `json:"amount"`
	PaymentMethod        string  `json:"payment_method"`
	Status               string  `json:"status"`
	CreatedAt            string  `json:"created_at"`
}

func (t *TransactionModel) CreateTransaction(transaction *Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO transactions (user_id, transaction_type, reference, amount, payment_method, status) VALUES (?,?,?,?,?,?)`

	if _, err := t.DB.ExecContext(ctx, query, transaction.UserID, transaction.TransactionType,
		transaction.TransactionReference, transaction.Amount,
		transaction.PaymentMethod, transaction.Status); err != nil {
		return err
	}

	return nil
}

func (t *TransactionModel) GetTransactions(userID uint) ([]Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, transaction_type, reference, amount, payment_method, status, created_at FROM transactions WHERE user_id = ?`

	rows, err := t.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var t Transaction

		if err := rows.Scan(
			&t.ID, 
			&t.UserID, 
			&t.TransactionType,
			&t.TransactionReference,
			&t.Amount,
			&t.PaymentMethod,
			&t.Status, 
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
