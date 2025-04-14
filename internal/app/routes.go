package app

import (
	"github.com/amir333/news-api-service/internal/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *Handler) {
	// Public routes
	public := r.Group("/api")
	{
		public.POST("/register", h.Register)
		public.POST("/login", h.Login)
		public.GET("/articles", h.GetArticles)
		public.GET("/articles/:id", h.GetArticle)
	}

	// Authenticated routes
	authenticated := r.Group("/api")
	authenticated.Use(auth.AuthMiddleware())
	{
		authenticated.POST("/articles", h.CreateArticle)
	}

	// Admin routes
	admin := r.Group("/api/admin")
	admin.Use(auth.AuthMiddleware(), auth.AdminMiddleware())
	{
		admin.GET("/users", h.GetAllUsers)
		admin.PUT("/articles/:id/status", h.UpdateArticleStatus)
	}
}
