package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"task-management/configs"
	"task-management/internal/db"
	"task-management/internal/handlers"
	"task-management/internal/middleware"
)

func main() {
	config := configs.LoadConfig()

	mongodb, err := db.NewMongoDB(config.MongoURI, config.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err := mongodb.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	jwtMiddleware := middleware.NewJwtMiddleware(config.JWTSecret)

	router := gin.Default()

	// Add CORS middleware to middleware chain
	router.Use(func(c *gin.Context) {
		log.Printf("Handling request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			log.Println("Handling OPTIONS request")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	handlers.SetupRoutes(router, mongodb.DB, jwtMiddleware)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Shutdown server with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited properly")
}
