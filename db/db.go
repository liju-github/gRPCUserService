package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/liju-github/EcommerceUserService/configs"
	"github.com/liju-github/EcommerceUserService/models"
)

// ConnectDB establishes a connection to the SQLite database using GORM
func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	// Use SQLite driver for GORM, and the database file will be located at the path specified in the config
	dsn := fmt.Sprintf("%s", "./db.sqlite3") // Example: use the file path from the config if available

	// Open the database connection using GORM
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool settings (typically used for relational databases)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting SQLDB instance from GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	// Auto-migrate the schema for the models (automatically creates or updates tables)
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("error auto-migrating database: %w", err)
	}

	log.Println("Successfully connected to SQLite database and auto-migrated the schema")
	return db, nil
}

// CloseDB safely closes the SQLite database connection
func CloseDB(db *gorm.DB) {
	// Ensure db is not nil before attempting to close the connection
	if db == nil {
		log.Println("No database connection to close.")
		return
	}

	// Get the underlying SQL database connection for closing
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting SQLDB instance: %v", err)
		return
	}

	// Attempt to close the database connection
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
		return
	}
	log.Println("SQLite database connection closed successfully")
}
