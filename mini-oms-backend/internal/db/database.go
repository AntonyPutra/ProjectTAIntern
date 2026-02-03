package db

import (
	"log"
	"mini-oms-backend/internal/config"
	"mini-oms-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes database connection
func Connect(cfg *config.Config) error {
	var err error

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if !cfg.IsDevelopment() {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// Connect to database
	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return err
	}

	log.Println("Database connected successfully")

	// Auto migrate models
	if err := AutoMigrate(); err != nil {
		return err
	}

	return nil
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate() error {
	log.Println("Running auto-migration...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.AuditLog{},
	)

	if err != nil {
		return err
	}

	log.Println("Auto-migration completed successfully")
	return nil
}

// GetDB returns database instance
func GetDB() *gorm.DB {
	return DB
}
