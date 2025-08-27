package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/malachi190/paylode-backend/models"
)

func AuthMiddleware(userModel models.UserModel) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenStr == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Bearer token is required",
			})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			ctx.Abort()
			return
		}

		// get user
		userId := claims["user_id"].(float64)

		user, err := userModel.GetUserById(uint(userId))

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unathorized access",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}