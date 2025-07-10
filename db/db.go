package db

import (
	"log"
)

// Init initializes the database connection and runs migrations
// This function maintains backward compatibility with existing code
func Init() {
	// Get database instance (this will trigger connection)
	_ = GetDB()

	// Run migrations
	if err := MigrateAll(); err != nil {
		log.Printf("Migration failed: %v", err)
		panic("failed to migrate database")
	}

	log.Println("Database initialized successfully")
}
