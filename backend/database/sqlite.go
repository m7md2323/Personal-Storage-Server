package database

import (
    "Personal-Storage-Server/backend/models" 
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    "log"
    "os"
    "github.com/joho/godotenv"
)

var DB *gorm.DB

func ConnectDatabase() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found or error reading it, using defaults if applicable")
    }

    path := os.Getenv("DATABASE_FILE_PATH")
    if path == "" {
        path = "./data_link/storage.db"
    }

    database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database!", err)
    }

    database.AutoMigrate(&models.File{})
	database.AutoMigrate(&models.User{})
    database.AutoMigrate(&models.Device{})

    DB = database
}