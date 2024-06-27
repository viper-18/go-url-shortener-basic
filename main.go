package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/viper-18/go-url-shortener-basic/models"
	"github.com/viper-18/go-url-shortener-basic/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func addHandlers(app *fiber.App) {
	app.Post("/shorten", routes.ShortenURL)
	app.Get("/:url", routes.NormaliseURL)
}

var DbCon *gorm.DB

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection parameters from environment variables
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create the DSN (Data Source Name) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUsername, dbPassword, dbHost, dbPort, dbName)

	// Open a connection to the database
	DbCon, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Ping the database to check if connection is successful
	sqlDB, err := DbCon.DB()
	if err != nil {
		log.Fatalf("Error getting DB: %v", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}

	log.Println("Connected to database successfully!")

	models.ConnectDB(DbCon)
	models.MigrateSchema()

	app := fiber.New()
	addHandlers(app)
	app.Listen(":80")
}
