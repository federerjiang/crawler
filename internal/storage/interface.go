package storage

import "time"

// Task represents a crawl task
type Task struct {
	ID          string
	StartURL    string
	Domain      string
	Status      string // pending, in_progress, completed, failed
	Links       []string
	TotalLinks  int
	ElapsedTime time.Duration
	CreatedAt   time.Time
	CompletedAt time.Time
}

// Storage defines the interface for task storage
type Storage interface {
	// SaveTask saves a new task
	SaveTask(task *Task) error

	// GetTask retrieves a task by its ID
	GetTask(taskID string) (*Task, error)

	// UpdateTaskStatus updates the status of a task
	UpdateTaskStatus(taskID, status string) error

	// AddTaskLinks adds links to a task
	AddTaskLinks(taskID string, links []string) error

	// SetTaskCompleted marks a task as completed with the elapsed time
	SetTaskCompleted(taskID string, elapsedTime time.Duration) error

	// SetTaskFailed marks a task as failed
	SetTaskFailed(taskID string, err error) error
}
