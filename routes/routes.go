package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/malachi190/paylode-backend/config"
	"github.com/malachi190/paylode-backend/handlers"
	"github.com/malachi190/paylode-backend/middleware"
)

func Router(d *config.Deps) *gin.Engine {
	r := gin.Default()

	// PLACE ROUTES HERE
	auth := r.Group("/api/auth")

	{
		auth.POST("/send-otp", handlers.SendOtpToken(d))
		auth.POST("/verify-email", handlers.VerifyEmail(d))
		auth.POST("/register", handlers.Register(d))
		auth.POST("/login", handlers.Login(d))
		auth.POST("/refresh", handlers.Refresh(d))
	}

	api := r.Group("/api").Use(middleware.AuthMiddleware(d.Models.Users))

	{
		api.POST("/create-pin", handlers.CreatePin(d))
		api.POST("/add-card", handlers.AddCard(d))
		api.POST("/fund-wallet/card", handlers.FundWalletWithCard(d))
		api.GET("/get-cards", handlers.FetchCards(d))
		api.GET("/get-wallet", handlers.GetUserWallet(d))
		api.GET("/fetch-transactions", handlers.FetchUserTransactions(d))
	}

	return r
}
