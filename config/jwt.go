package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateRefreshToken(userID any) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}).SignedString([]byte(os.Getenv("RefreshJWTSecret")))
}

func GenerateAuthToken(userID any) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (d *Deps) ValidateRefreshToken(refreshToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported signing method")
		}
		return []byte(os.Getenv("RefreshJWTSecret")), nil
	})

	if err != nil || !token.Valid {
		ErrorLogger.Printf("invalid refresh token: %v\n", refreshToken)

		return nil, fmt.Errorf("invalid token provided")
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		ErrorLogger.Printf("invalid token type: %v\n", claims["type"])

		return nil, fmt.Errorf("invalid token provided")
	}

	// CHECK DB REVOCATION
	hash := sha256.Sum256([]byte(refreshToken))
	hashEncode := hex.EncodeToString(hash[:])

	ok, err := d.Models.Sessions.RefreshTokenExists(hashEncode)

	if !ok || err != nil {
		return nil, fmt.Errorf("token revoked or missing")
	}

	return claims, nil
}
