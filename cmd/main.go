package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/malachi190/paylode-backend/config"
	"github.com/malachi190/paylode-backend/models"
	"github.com/malachi190/paylode-backend/routes"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// DECLARE MODELS
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	defer db.Close()

	models := models.HandleModels(db)

	// INITIALIZE REDIS
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})

	// SET OPTIONS IN DEPENDENCY
	deps := &config.Deps{
		Models: models,
		Redis:  rdb,
	}

	// TODO: Implement rate limiting using redis

	// SET UP ROUTES
	router := routes.Router(deps)

	// SET UP SERVER OPTIONS
	port := os.Getenv("PORT")

	addr := "0.0.0.0:" + port

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
