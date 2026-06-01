package model

import (
	"errors"
	"strings"
	"time"
)

// Todo represents a discrete unit of work.
type Todo struct {
	ID          int                `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Completed   bool               `json:"completed"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CompletedAt *time.Time         `json:"completed_at"`
	DueDate     string             `json:"due_date"`
}

// Validate checks that the Todo has a valid title and due_date format.
func (t Todo) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title is required")
	}

	if t.DueDate != "" {
		if _, err := time.Parse("2006-01-02", t.DueDate); err != nil {
			return errors.New("due_date must be in YYYY-MM-DD format")
		}
	}

	return nil
}
