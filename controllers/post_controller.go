package controllers

import (
	"blog-api/config"
	"blog-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// GET /api/posts - Lấy tất cả bài viết
func FindPosts(c *gin.Context) {
	var posts []models.Post
	// Preload("Author") để lấy luôn thông tin tác giả (JOIN bảng)
	config.DB.Preload("Author").Find(&posts)
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// GET /api/posts/:id - Lấy chi tiết 1 bài viết
func FindPost(c *gin.Context) {
	var post models.Post
	if err := config.DB.Preload("Author").Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bài viết!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}

// POST /api/posts - Tạo mới bài viết (Cần Token)
func CreatePost(c *gin.Context) {
	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lấy userID từ middleware xác thực truyền qua
	userID, _ := c.Get("userID")

	post := models.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID.(uint),
	}

	config.DB.Create(&post)
	c.JSON(http.StatusCreated, gin.H{"data": post})
}

// PUT /api/posts/:id - Cập nhật bài viết (Chỉ tác giả)
func UpdatePost(c *gin.Context) {
	var post models.Post
	if err := config.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bài viết!"})
		return
	}

	userID, _ := c.Get("userID")
	if post.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền chỉnh sửa bài viết của người khác!"})
		return
	}

	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&post).Updates(models.Post{Title: input.Title, Content: input.Content})
	c.JSON(http.StatusOK, gin.H{"data": post})
}

// DELETE /api/posts/:id - Xóa bài viết (Chỉ tác giả)
func DeletePost(c *gin.Context) {
	var post models.Post
	if err := config.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bài viết!"})
		return
	}

	userID, _ := c.Get("userID")
	if post.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền xóa bài viết của người khác!"})
		return
	}

	config.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Xóa bài viết thành công!"})
}
