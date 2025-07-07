package db

import (
	"fmt"
	"log"
	"os"
	"sykell-challenge/backend/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	// Get database configuration from environment variables with defaults
	dbUser := getEnv("DB_USER", "sykell")
	dbPassword := getEnv("DB_PASSWORD", "sykellpass")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "websites_dev")

	// Construct DSN using environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Retry logic for database connection
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Printf("Successfully connected to database on attempt %d", i+1)
			return
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(2 * time.Second)
		}
	}

	panic(fmt.Sprintf("failed to connect database after %d attempts: %v", maxRetries, err))
}
func GetDB() *gorm.DB {
	if DB == nil {
		Init()
	}
	return DB
}
func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		panic("failed to close database connection")
	}
	err = sqlDB.Close()
	if err != nil {
		panic("failed to close database connection")
	}
}
func Migrate() {
	if DB == nil {
		Init()
	}
	err = DB.AutoMigrate(
		&models.URL{},
		&models.User{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
