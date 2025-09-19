package routes

import (
	"rentx/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// authentication routes
	server.POST("/refresh-token", refreshTokenHandler)
	server.POST("/auth-email", emailAuthHandler)
	server.POST("/auth-phone", phoneAuthHandler)
	server.POST("/auth-oauth", oauthAuthHandler)
	// server.POST("/signup", signup)
	// server.POST("/login", login)
	// categories
	server.POST("/category", middlewares.Authenticate, createCategory)
	server.PUT("/category/:id", middlewares.Authenticate, updateCategory)
	server.DELETE("/category/:id", middlewares.Authenticate, deleteCategory)
	server.GET("/category", listCategories)
	// posts
	server.POST("/upload/post-image", middlewares.Authenticate, uploadPostImages)
	server.POST("/posts", middlewares.Authenticate, createPost)
	server.PUT("/posts/:id", middlewares.Authenticate, updatePost)
	server.DELETE("/posts/:id", middlewares.Authenticate, deletePost)
	server.GET("/posts", listPosts)
	server.GET("/posts/:id", getPostByID)
	// orders
	server.POST("/orders", createOrder)
	server.DELETE("/orders/:id", deleteOrder)
	server.GET("/orders", listOrders)
	server.GET("/orders/:id", getOrderByID)
	// reviews
	server.POST("/reviews", middlewares.Authenticate, createReview)
	server.GET("/reviews/:postId", listReviewsByPost)
	server.DELETE("/reviews/:id", middlewares.Authenticate, deleteReview)
}
