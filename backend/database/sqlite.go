package database

import (
    "Personal-Storage-Server/backend/models" 
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

func ConnectDatabase() {
    path := "./data_link/storage.db"

    database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database!", err)
    }

    database.AutoMigrate(&models.File{})
	database.AutoMigrate(&models.User{})
    database.AutoMigrate(&models.Device{})

    DB = database
}