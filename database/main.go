package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/golang-migrate/migrate/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please indicate migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]

	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
		log.Print(err.Error())
	}

	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Build MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	fmt.Printf("DSN built: %s\n", dsn)

	// open database
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Panicf("error initializing database connection: %s\n", err.Error())
	}

	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Ping to make sure it works
	if err := db.Ping(); err != nil {
		log.Fatalf("ping: %v", err)
	}

	// Create migration driver instance
	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		log.Fatalf("error loading migrate driver: %v", err)
	}

	// Get source file
	file, err := (&file.File{}).Open("./migrations")

	if err != nil {
		log.Fatalf("failed to load file: %v", err)
	}

	// Create source and migrate instance
	m, err := migrate.NewWithInstance("file", file, "mysql", driver)

	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// Apply migrations
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	case "down":
		// roll back only the last migration
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'.")
	}

}
