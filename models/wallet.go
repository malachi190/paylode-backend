package models

import (
	"context"
	"database/sql"
	"time"
)

type WalletModel struct {
	DB *sql.DB
}

type Wallet struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	WalletBalance float64 `json:"wallet_balance"`
	WalletID      string  `json:"wallet_id"`
	Currency      string  `json:"currency"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	User          *User   `json:"user,omitempty"`
}

func (w *WalletModel) CreateWallet(wallet *Wallet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO wallets (user_id, wallet_balance, wallet_id) VALUES (?,?,?)`

	_, err := w.DB.ExecContext(ctx, query, wallet.UserID, wallet.WalletBalance, wallet.WalletID)

	if err != nil {
		return err
	}

	return nil
}

func (w *WalletModel) Fund(userID uint, amount float64) (*Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	amountInKobo := int64(amount * 100)

	// Begin db transaction
	tx, _ := w.DB.BeginTx(ctx, nil)
	defer tx.Rollback()

	var wallet Wallet

	if err := w.DB.QueryRowContext(ctx, `SELECT id, user_id, wallet_balance, wallet_id, currency, created_at, updated_at FROM wallets WHERE user_id = ? FOR UPDATE`, userID).
		Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.WalletBalance,
			&wallet.WalletID,
			&wallet.Currency,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
		); err != nil {
		return nil, err
	}

	newBal := wallet.WalletBalance + float64(amountInKobo)

	if _, err := w.DB.ExecContext(ctx, `UPDATE wallets SET wallet_balance = ? WHERE id = ?`, newBal, wallet.ID); err != nil {
		return nil, err
	}

	wallet.WalletBalance = newBal

	return &wallet, tx.Commit()

}

func (w *WalletModel) GetWallet(userID uint) (*Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, wallet_balance, wallet_id, currency, created_at, updated_at FROM wallets WHERE user_id = ?`

	var wallet Wallet

	err := w.DB.QueryRowContext(ctx, query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.WalletBalance,
		&wallet.WalletID,
		&wallet.Currency,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
