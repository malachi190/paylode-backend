package config

import (
	"github.com/malachi190/paylode-backend/models"
	"github.com/redis/go-redis/v9"
)



type Deps struct {
	Models models.Models
	Redis *redis.Client
}