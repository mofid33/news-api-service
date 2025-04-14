package app

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:"default:false"`
}

type Article struct {
	gorm.Model
	Title       string    `gorm:"not null"`
	Content     string    `gorm:"type:text;not null"`
	AuthorID    uint      `gorm:"not null"`
	Author      User      `gorm:"foreignKey:AuthorID"`
	PublishedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status      string    `gorm:"default:'draft'"` // draft, published, archived
	Category    string
	Tags        []Tag `gorm:"many2many:article_tags;"`
}

type Tag struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}