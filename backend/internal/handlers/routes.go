package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"task-management/internal/middleware"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(router *gin.Engine, db *mongo.Database, jwtMiddleware *middleware.JwtMiddleware) {
	// Collections
	userCollection := db.Collection("users")
	goalCollection := db.Collection("goals")

	// Handlers
	authHandler := NewAuthHandler(userCollection, jwtMiddleware)
	goalHandler := NewGoalHandler(goalCollection)

	// Auth routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Goal routes (protected)
	goals := router.Group("/api/goals")
	goals.Use(jwtMiddleware.AuthRequired())
	{
		goals.POST("", goalHandler.CreateGoal)
		goals.GET("", goalHandler.ListGoals)
		goals.GET("/:id", goalHandler.GetGoal)
		goals.PUT("/:id", goalHandler.UpdateGoal)
		goals.DELETE("/:id", goalHandler.DeleteGoal)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
