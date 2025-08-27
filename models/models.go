package models

import "database/sql"

type Models struct {
	Users    UserModel
	Sessions SessionModel
}

func HandleModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Sessions: SessionModel{DB: db},
	}
}
