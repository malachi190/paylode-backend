package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/malachi190/paylode-backend/config"
	"github.com/malachi190/paylode-backend/handlers"
)

func Router(d *config.Deps) *gin.Engine {
	r := gin.Default()

	// PLACE ROUTES HERE
	api := r.Group("/api/auth")

	{
		api.POST("/send-otp", handlers.SendOtpToken(d))
		api.POST("/verify-email", handlers.VerifyEmail(d))
		api.POST("/register", handlers.Register(d))
		api.POST("/login", handlers.Login(d))
		api.POST("/refresh", handlers.Refresh(d))
	}

	return r
}
