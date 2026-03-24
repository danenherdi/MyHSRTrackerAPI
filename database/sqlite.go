package database

import (
	"log"
	"os"

	"MyHSRTrackerAPI/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	
	dbName := os.Getenv("DB_NAME")

	DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.WarpLog{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	log.Println("Database connected and migrated successfully.")
}
