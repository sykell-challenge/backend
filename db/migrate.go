package db

import (
	"sykell-challenge/backend/models"
)

// MigrateAll runs auto-migration for all models
func MigrateAll() error {
	db := GetDB()
	return db.AutoMigrate(
		&models.URL{},
		&models.User{},
		&models.CrawlJob{},
	)
}
