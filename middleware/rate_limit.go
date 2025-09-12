package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func RateLimit(limiter *redis_rate.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := limiter.Allow(ctx, ctx.ClientIP(), redis_rate.PerMinute(10))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}

		if res.Remaining == 0 {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests",
			})
			return
		}

		ctx.Next()
	}
}
