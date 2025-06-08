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
	"go.mongodb.org/mongo-driver/mongo/options"

	"task-management/internal/models"
)

// GoalHandler handles goal related routes
type GoalHandler struct {
	goalCollection *mongo.Collection
	validator      *validator.Validate
}

// NewGoalHandler creates a new goal handler
func NewGoalHandler(goalCollection *mongo.Collection) *GoalHandler {
	return &GoalHandler{
		goalCollection: goalCollection,
		validator:      validator.New(),
	}
}

// CreateGoalRequest represents the create goal request
type CreateGoalRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description,omitempty"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
}

// UpdateGoalRequest represents the update goal request
type UpdateGoalRequest struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	StartDate   time.Time  `json:"startDate,omitempty"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Completed   bool       `json:"completed,omitempty"`
}

// AddSubTaskRequest represents the add subtask request
type AddSubTaskRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

// CreateGoal handles goal creation
func (h *GoalHandler) CreateGoal(c *gin.Context) {
	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	now := time.Now()
	goal := models.Goal{
		ID:          primitive.NewObjectID(),
		UserID:      userID.(primitive.ObjectID),
		Title:       req.Title,
		Description: req.Description,
		SubTasks:    []models.SubTask{},
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Completed:   false,
		Progress:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Insert goal to database
	_, err := h.goalCollection.InsertOne(context.Background(), goal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// GetGoal handles getting a single goal
func (h *GoalHandler) GetGoal(c *gin.Context) {
	goalID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Find goal by ID and user ID
	var goal models.Goal
	err = h.goalCollection.FindOne(context.Background(), bson.M{
		"_id":    goalID,
		"userId": userID.(primitive.ObjectID),
	}).Decode(&goal)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get goal"})
		}
		return
	}

	c.JSON(http.StatusOK, goal)
}

// ListGoals handles listing all goals for a user
func (h *GoalHandler) ListGoals(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Find all goals by user ID
	cursor, err := h.goalCollection.Find(context.Background(),
		bson.M{"userId": userID.(primitive.ObjectID)},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list goals"})
		return
	}
	defer cursor.Close(context.Background())

	var goals []models.Goal
	if err := cursor.All(context.Background(), &goals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode goals"})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// UpdateGoal handles updating a goal
func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	goalID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var req UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"updatedAt": time.Now(),
	}

	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if !req.StartDate.IsZero() {
		update["startDate"] = req.StartDate
	}
	if req.EndDate != nil {
		update["endDate"] = req.EndDate
	}
	update["completed"] = req.Completed

	result, err := h.goalCollection.UpdateOne(
		context.Background(),
		bson.M{
			"_id":    goalID,
			"userId": userID.(primitive.ObjectID),
		},
		bson.M{"$set": update},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update goal"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Get updated goal
	var goal models.Goal
	err = h.goalCollection.FindOne(context.Background(), bson.M{
		"_id":    goalID,
		"userId": userID.(primitive.ObjectID),
	}).Decode(&goal)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated goal"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteGoal handles deleting a goal
func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	goalID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	result, err := h.goalCollection.DeleteOne(context.Background(), bson.M{
		"_id":    goalID,
		"userId": userID.(primitive.ObjectID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete goal"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}
