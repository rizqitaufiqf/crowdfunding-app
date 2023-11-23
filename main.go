package main

import (
	"crowdfunding/auth"
	"crowdfunding/handler"
	"crowdfunding/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	// load env file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	// setup env values to variable
	host := os.Getenv("DATABASE_HOST")
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	port := os.Getenv("DATABASE_PORT")
	sslMode := os.Getenv("DATABASE_SSL_ENABLED")
	timezone := os.Getenv("TIMEZONE")

	// generate connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, username, password, databaseName, port, sslMode, timezone)
	// connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection to database successfully")
	// init User Repository
	userRepository := user.NewRepository(db)
	// init User Service
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	// init User Handler(Controller)
	userHandler := handler.NewUserHandler(userService, authService)

	// init router
	router := gin.Default()
	// group router endpoint
	api := router.Group("/api/v1")
	// generate endpoint
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/check-email", userHandler.CheckEmail)
	api.POST("/avatars", userHandler.UploadAvatar)
	// run web service
	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatal(err.Error())
	}

}
