package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("DATABASE_HOST")
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	port := os.Getenv("DATABASE_PORT")
	sslMode := os.Getenv("DATABASE_SSL_ENABLED")
	timezone := os.Getenv("TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, username, password, databaseName, port, sslMode, timezone)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection to database successfully")

}
