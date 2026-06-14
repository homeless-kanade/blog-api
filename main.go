package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Blog struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

var blogs = []Blog{
	{
		ID:      1,
		Title:   "Learn Golang",
		Content: "Golang is awesome",
		Author:  "Lam",
	},
}

func main() {
	r := gin.Default()

	r.GET("/blogs", getBlogs)
	r.GET("/blogs/:id", getBlogByID)
	r.POST("/blogs", createBlog)
	r.PUT("/blogs/:id", updateBlog)
	r.DELETE("/blogs/:id", deleteBlog)

	r.Run(":8080")
}

func getBlogs(c *gin.Context) {
	c.JSON(http.StatusOK, blogs)
}

func getBlogByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for _, blog := range blogs {
		if blog.ID == id {
			c.JSON(http.StatusOK, blog)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
}

func createBlog(c *gin.Context) {
	var newBlog Blog

	if err := c.ShouldBindJSON(&newBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBlog.ID = len(blogs) + 1
	blogs = append(blogs, newBlog)

	c.JSON(http.StatusCreated, newBlog)
}

func updateBlog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedBlog Blog
	if err := c.ShouldBindJSON(&updatedBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, blog := range blogs {
		if blog.ID == id {
			updatedBlog.ID = id
			blogs[i] = updatedBlog
			c.JSON(http.StatusOK, updatedBlog)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
}

func deleteBlog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for i, blog := range blogs {
		if blog.ID == id {
			blogs = append(blogs[:i], blogs[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
}
