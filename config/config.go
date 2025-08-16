package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	// Only try to load .env file in development
	// In production, environment variables should be set directly
	if os.Getenv("ENVIRONMENT") == "" || os.Getenv("ENVIRONMENT") == "development" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Printf("Warning: .env file not found (this is normal in production): %v", err)
		} else {
			log.Println("Loaded .env file successfully")
		}
	} else {
		log.Println("Running in production mode - using environment variables directly")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func ConnectDB() {
	var err error
	dsn := GetEnv("DB_URI")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Error connecting to DB:%s", err)
	}

	// err = DB.AutoMigrate(&models.User{}, &models.Album{}, &models.Song{}, &models.Playlist{})
	if err != nil {
		log.Fatalf("Error migrating models:%s", err)
	}

	log.Println("Connected to DB successfully 🎉🎉")
}
