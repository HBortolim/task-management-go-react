package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"task-management/internal/middleware"
	"task-management/internal/models"
	"task-management/internal/utils"
)

// AuthHandler handles authentication related routes
type AuthHandler struct {
	userCollection *mongo.Collection
	validator      *validator.Validate
	jwtMiddleware  *middleware.JwtMiddleware
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userCollection *mongo.Collection, jwtMiddleware *middleware.JwtMiddleware) *AuthHandler {
	return &AuthHandler{
		userCollection: userCollection,
		validator:      validator.New(),
		jwtMiddleware:  jwtMiddleware,
	}
}

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// LoginRequest represents the user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user with email already exists
	var existingUser models.User
	err := h.userCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Check if username is taken
	err = h.userCollection.FindOne(context.Background(), bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
		return
	}

	// Create new user
	now := time.Now()
	user := models.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user to database
	_, err = h.userCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := h.jwtMiddleware.GenerateToken(user.ID, 24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return token and user
	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	defer utils.TraceFunction()()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug login request
	utils.DebugObject("Login Request", req)

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	err := h.userCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare passwords
	if err := user.ComparePassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := h.jwtMiddleware.GenerateToken(user.ID, 24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return token and user
	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}
