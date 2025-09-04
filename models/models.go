package models

import "database/sql"

type Models struct {
	Users        UserModel
	Sessions     SessionModel
	Wallets      WalletModel
	Cards        CardModel
	Transactions TransactionModel
}

func HandleModels(db *sql.DB) Models {
	return Models{
		Users:        UserModel{DB: db},
		Sessions:     SessionModel{DB: db},
		Wallets:      WalletModel{DB: db},
		Cards:        CardModel{DB: db},
		Transactions: TransactionModel{DB: db},
	}
}
