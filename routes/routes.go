package routes

import (
	"blog-api/controllers"
	"blog-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Các router public không cần Token
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
	}

	postRoutes := r.Group("/api/posts")
	{
		postRoutes.GET("", controllers.FindPosts)
		postRoutes.GET("/:id", controllers.FindPost)

		// Nhóm API cần phải qua lớp kiểm tra Token
		protected := postRoutes.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("", controllers.CreatePost)
			protected.PUT("/:id", controllers.UpdatePost)
			protected.DELETE("/:id", controllers.DeletePost)
		}
	}

	return r
}