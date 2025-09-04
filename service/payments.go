package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type ApiGateway struct{}

// type AddMoneyRequest struct {
// 	Amount   int64
// 	Currency string
// 	Token    string
// }

// type AddMoneyResult struct {
// 	Success bool
// 	ID      string // transaction id
// 	Message string
// }

func AddCard(card_number, brand string, exp_month, exp_year, cvv string) (token string, last_four string, err error) {
	if len(card_number) < 4 {
		return "", "", errors.New("invalid card number")
	}

	// GENERATE MOCK TOKEN --- WHEN USING LIVE GATEWAY SERVICE, TOKEN WOULD BE GENERATED FROM THERE AND RETURNED
	b := make([]byte, 16)
	rand.Read(b)
	token = hex.EncodeToString(b)
	return token, card_number[len(card_number)-4:], nil
}

func ChargeCard(token string) (transactionRef string, err error) {
	if time.Now().UnixNano()%20 == 0 {
		return "", errors.New("insufficient funds")
	}

	return "txn_" + token[:8], nil
}
