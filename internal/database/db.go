package database

import (
	"fmt"
	"log/slog"

	"jellyfin-telegram-bot/pkg/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the gorm.DB instance
type DB struct {
	*gorm.DB
}

// NewDB creates a new database connection
func NewDB(dbPath string) (*DB, error) {
	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Use Silent in production, Info for debugging
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	slog.Info("Connected to database", "path", dbPath)

	// Auto-migrate schema
	if err := db.AutoMigrate(&models.Subscriber{}, &models.ContentCache{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	slog.Info("Database schema migrated successfully")

	return &DB{DB: db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
