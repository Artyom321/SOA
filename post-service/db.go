package main

import (
	"log"
	"time"

	"social-network/common/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDB(connStr string) *gorm.DB {
	var db *gorm.DB
	var err error

	maxRetries := 10
	retryDelay := time.Second * 3

	for retries := 0; retries < maxRetries; retries++ {
		log.Printf("Attempting to connect to database %d/%d", retries+1, maxRetries)
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to database")
			return db
		}

		log.Printf("Database connection error: %v", err)
		log.Printf("Retrying in %v", retryDelay)
		time.Sleep(retryDelay)
	}

	log.Fatalf("Failed to connect to database after %d attempts", maxRetries)
	return nil
}

func initDB() {
	err := db.AutoMigrate(&models.Post{}, &models.Comment{}, &models.Like{})
	if err != nil {
		log.Fatalf("Database migration error: %v", err)
	}
	log.Println("Database migration completed successfully")
}
