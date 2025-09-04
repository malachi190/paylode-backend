package config

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/malachi190/paylode-backend/models"
)

func GenerateOtp() int {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return int(n.Int64() + 100000)
}

func GenerateWalletID() int {
	n, _ := rand.Int(rand.Reader, big.NewInt(9000000000))
	return int(n.Int64() + 10000000000)
}

func GetLoggedInUser(ctx *gin.Context) (*models.User, error) {
	raw, exists := ctx.Get("user")

	if !exists {
		return nil, errors.New("unauthorized access")
	}

	user, ok := raw.(*models.User)

	if !ok {
		return nil, errors.New("invalid user type")
	}

	return user, nil
}
