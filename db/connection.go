package db

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"sykell-challenge/backend/utils"
)

var (
	instance *gorm.DB
	once     sync.Once
	initErr  error
)

// Config holds database configuration
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// LoadConfig loads database configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		User:     utils.GetEnv("DB_USER", "sykell"),
		Password: utils.GetEnv("DB_PASSWORD", "sykellpass"),
		Host:     utils.GetEnv("DB_HOST", "127.0.0.1"),
		Port:     utils.GetEnv("DB_PORT", "3306"),
		Name:     utils.GetEnv("DB_NAME", "websites_dev"),
	}
}

// Connect establishes a connection to the database
func Connect(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Name)

	// Retry logic for database connection
	const maxRetries = 30
	const retryDelay = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
		if err == nil {
			log.Printf("Successfully connected to database on attempt %d", i+1)
			return db, nil
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect database after %d attempts", maxRetries)
}

// GetDB returns the database instance, initializing it if necessary
func GetDB() *gorm.DB {
	once.Do(func() {
		config := LoadConfig()
		instance, initErr = Connect(config)
	})

	if initErr != nil {
		panic(fmt.Sprintf("failed to initialize database: %v", initErr))
	}

	return instance
}
