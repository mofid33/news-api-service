package main

import (
	"log"

	"gorm.io/gorm"

	"github.com/amir333/news-api-service/internal/app"
	"github.com/amir333/news-api-service/internal/auth"
	"github.com/amir333/news-api-service/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	dsn := "host=localhost user=news_api_user password=Aa@123456 dbname=news_api port=5432 sslmode=disable"
	db, err := storage.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create admin user if not exists
	createAdminUser(db)

	// Initialize Gin
	r := gin.Default()

	// Set up routes
	handler := &app.Handler{DB: db}
	app.SetupRoutes(r, handler)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&app.User{}).Where("is_admin = ?", true).Count(&count)

	if count == 0 {
		hashedPassword, _ := auth.HashPassword("admin123")
		admin := app.User{
			Username: "admin",
			Email:    "admin@newsapi.com",
			Password: hashedPassword,
			IsAdmin:  true,
		}
		db.Create(&admin)
	}
}
