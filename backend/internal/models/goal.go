package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubTask represents a subtask within a goal
type SubTask struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Completed   bool               `json:"completed" bson:"completed"`
	DueDate     *time.Time         `json:"dueDate,omitempty" bson:"dueDate,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// Goal represents a user's goal
type Goal struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	SubTasks    []SubTask          `json:"subTasks" bson:"subTasks"`
	StartDate   time.Time          `json:"startDate" bson:"startDate"`
	EndDate     *time.Time         `json:"endDate,omitempty" bson:"endDate,omitempty"`
	Completed   bool               `json:"completed" bson:"completed"`
	Progress    float64            `json:"progress" bson:"progress"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// CalculateProgress calculates the progress of a goal based on completed subtasks
func (g *Goal) CalculateProgress() {
	if len(g.SubTasks) == 0 {
		g.Progress = 0
		return
	}

	completedCount := 0
	for _, task := range g.SubTasks {
		if task.Completed {
			completedCount++
		}
	}

	g.Progress = float64(completedCount) / float64(len(g.SubTasks)) * 100
}

// IsCompleted checks if all subtasks are completed and updates the goal status
func (g *Goal) IsCompleted() bool {
	if len(g.SubTasks) == 0 {
		return false
	}

	for _, task := range g.SubTasks {
		if !task.Completed {
			g.Completed = false
			return false
		}
	}

	g.Completed = true
	return true
}
