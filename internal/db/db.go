package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"goCachedAPI/internal/models"
)

func New(sqliteDSN string) *gorm.DB {
	database, err := gorm.Open(sqlite.Open(sqliteDSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to sqlite: %v", err)
	}

	if err := database.AutoMigrate(&models.Product{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	seed(database)

	return database
}

func seed(database *gorm.DB) {
	var count int64

	if err := database.Model(&models.Product{}).Count(&count).Error; err != nil {
		log.Fatalf("failed to count products: %v", err)
	}

	if count == 0 {
		if err := database.Create(&models.Product{
			Name:  "Product1",
			Price: 100,
		}).Error; err != nil {
			log.Fatalf("failed to seed products: %v", err)
		}
	}
}
