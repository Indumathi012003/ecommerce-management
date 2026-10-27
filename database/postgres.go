package database

import (
	"fmt"
	"time"

	"order-management/config"
	"order-management/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Optional: AutoMigrate tables directly with GORM
	if err := db.AutoMigrate(&model.Order{}, &model.Product{}); err != nil {
		return nil, fmt.Errorf("auto-migration failed: %w", err)
	}

	// Configure underlying connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}