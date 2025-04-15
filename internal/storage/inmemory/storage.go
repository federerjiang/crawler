package inmemory

import (
	"errors"
	"sync"
	"time"

	"github.com/crawler/internal/storage"
)

// InMemoryStorage implements storage.Storage interface with in-memory storage
type InMemoryStorage struct {
	tasks map[string]*storage.Task
	mutex sync.RWMutex
}

// NewStorage creates a new instance of InMemoryStorage
func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: make(map[string]*storage.Task),
	}
}

// SaveTask saves a new task
func (s *InMemoryStorage) SaveTask(task *storage.Task) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.tasks[task.ID] = task
	return nil
}

// GetTask retrieves a task by its ID
func (s *InMemoryStorage) GetTask(taskID string) (*storage.Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

// UpdateTaskStatus updates the status of a task
func (s *InMemoryStorage) UpdateTaskStatus(taskID, status string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return errors.New("task not found")
	}

	task.Status = status
	return nil
}

// AddTaskLinks adds links to a task
func (s *InMemoryStorage) AddTaskLinks(taskID string, links []string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return errors.New("task not found")
	}

	// Add only unique links
	linksMap := make(map[string]bool)
	for _, link := range task.Links {
		linksMap[link] = true
	}

	for _, link := range links {
		if _, exists := linksMap[link]; !exists {
			task.Links = append(task.Links, link)
			linksMap[link] = true
		}
	}

	task.TotalLinks = len(task.Links)
	return nil
}

// SetTaskCompleted marks a task as completed with the elapsed time
func (s *InMemoryStorage) SetTaskCompleted(taskID string, elapsedTime time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return errors.New("task not found")
	}

	task.Status = "completed"
	task.ElapsedTime = elapsedTime
	task.CompletedAt = time.Now()
	return nil
}

// SetTaskFailed marks a task as failed
func (s *InMemoryStorage) SetTaskFailed(taskID string, err error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return errors.New("task not found")
	}

	task.Status = "failed"
	task.CompletedAt = time.Now()
	return nil
}
