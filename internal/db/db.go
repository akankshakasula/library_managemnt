package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"library-management/internal/models"
)

// DB is the global GORM database connection instance
var DB *gorm.DB

// JWTSecret is the global JWT secret key
var JWTSecret string // ADD THIS LINE

// ConnectDatabase initializes the database connection and performs migrations
func ConnectDatabase() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	projectRoot := filepath.Join(basepath, "../../")
	envPath := filepath.Join(projectRoot, ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL not set in .env file")
	}

	// Load JWT Secret
	JWTSecret = os.Getenv("JWT_SECRET") // ADD THIS LINE
	if JWTSecret == "" {                // ADD THIS LINE
		log.Fatal("JWT_SECRET not set in .env file") // ADD THIS LINE
	} // ADD THIS LINE

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database using URL '%s': %v", databaseURL, err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Borrow{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	fmt.Println("Database Migrated: User, Book, and Borrow tables created/updated")

	DB = db
	fmt.Println("Connected to Database!")
}
