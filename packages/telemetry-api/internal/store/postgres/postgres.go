package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"omoidememo.com/telemetry-api/internal/telemetry"
)

func New(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&telemetry.Event{}); err != nil {
		return nil, err
	}
	return db, nil
}
