package app

import (
	"net/http"

	"github.com/amir333/news-api-service/internal/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

// User handlers
func (h *Handler) Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user.Password = hashedPassword
	result := h.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := h.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !auth.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Article handlers
func (h *Handler) CreateArticle(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var article Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article.AuthorID = userID
	result := h.DB.Create(&article)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (h *Handler) GetArticles(c *gin.Context) {
	var articles []Article
	if err := h.DB.Preload("Author").Preload("Tags").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (h *Handler) GetArticle(c *gin.Context) {
	id := c.Param("id")

	var article Article
	if err := h.DB.Preload("Author").Preload("Tags").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// Admin handlers
func (h *Handler) GetAllUsers(c *gin.Context) {
	var users []User
	if err := h.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove passwords from response
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateArticleStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var article Article
	if err := h.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	article.Status = input.Status
	h.DB.Save(&article)

	c.JSON(http.StatusOK, article)
}
