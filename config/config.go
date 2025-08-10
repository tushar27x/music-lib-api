package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tushar27x/music-lib-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file:%v", err)
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

	err = DB.AutoMigrate(&models.User{}, &models.Album{}, &models.Song{}, &models.Playlist{})
	if err != nil {
		log.Fatalf("Error migrating models:%s", err)
	}

	log.Println("Connected to DB successfully 🎉🎉")
}
