package storage

import (
	"fmt"

	"github.com/amir333/news-api-service/internal/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&app.User{}, &app.Article{}, &app.Tag{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
