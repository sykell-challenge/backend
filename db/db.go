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
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
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
		&models.CrawlJob{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
}

// Development utilities

// DropAllTables drops all tables (complete reset)
func DropAllTables() error {
	if DB == nil {
		Init()
	}

	// Drop tables in reverse order of dependencies
	tables := []interface{}{
		&models.CrawlJob{},
		&models.URL{},
		&models.User{},
	}

	for _, table := range tables {
		if err := DB.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table %T: %v", table, err)
		}
	}

	return nil
}

// ClearAllData deletes all records but keeps table structure
func ClearAllData() error {
	if DB == nil {
		Init()
	}

	// Disable foreign key checks (MySQL specific)
	if err := DB.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %v", err)
	}

	// Delete in order that respects foreign key constraints
	if err := DB.Where("1 = 1").Delete(&models.CrawlJob{}).Error; err != nil {
		return fmt.Errorf("failed to delete crawl jobs: %v", err)
	}

	if err := DB.Where("1 = 1").Delete(&models.URL{}).Error; err != nil {
		return fmt.Errorf("failed to delete URLs: %v", err)
	}

	if err := DB.Where("1 = 1").Delete(&models.User{}).Error; err != nil {
		return fmt.Errorf("failed to delete users: %v", err)
	}

	// Re-enable foreign key checks
	if err := DB.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("failed to re-enable foreign key checks: %v", err)
	}

	return nil
}

// GetTableCounts returns record counts for all tables
func GetTableCounts() map[string]int64 {
	if DB == nil {
		Init()
	}

	counts := make(map[string]int64)
	var count int64

	DB.Model(&models.URL{}).Count(&count)
	counts["urls"] = count

	DB.Model(&models.User{}).Count(&count)
	counts["users"] = count

	DB.Model(&models.CrawlJob{}).Count(&count)
	counts["crawl_jobs"] = count

	return counts
}

// ResetDatabase completely resets the database (drop + migrate)
func ResetDatabase() error {
	if err := DropAllTables(); err != nil {
		return fmt.Errorf("failed to drop tables: %v", err)
	}

	Migrate()
	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
