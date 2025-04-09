package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/julioCAlmeida/go-gateway/internal/repository"
	"github.com/julioCAlmeida/go-gateway/internal/service"
	"github.com/julioCAlmeida/go-gateway/internal/web/server"
	_ "github.com/lib/pq"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "postgres"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	port := getEnv("PORT", "8080")
	server := server.NewServer(accountService, port)
	server.ConfigureRoutes()
	
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	
	log.Println("Server started on port", port)
}