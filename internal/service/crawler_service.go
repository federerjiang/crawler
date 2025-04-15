package service

import (
	"fmt"
	"time"

	"github.com/crawler/internal/crawler"
	"github.com/crawler/internal/storage"
	"github.com/google/uuid"
)

// CrawlerService is responsible for managing crawling tasks
type CrawlerService struct {
	storage storage.Storage
}

// NewCrawlerService creates a new CrawlerService
func NewCrawlerService(storage storage.Storage) *CrawlerService {
	return &CrawlerService{
		storage: storage,
	}
}

// StartCrawl starts a new crawling task
func (s *CrawlerService) StartCrawl(startURL string, concurrentRequests, maxDepth int, timeout time.Duration) (string, error) {
	// Create task ID
	taskID := uuid.New().String()

	// Extract domain
	domain, err := crawler.ExtractDomain(startURL)
	if err != nil {
		return "", fmt.Errorf("invalid start URL: %v", err)
	}

	// Create task
	task := &storage.Task{
		ID:         taskID,
		StartURL:   startURL,
		Domain:     domain,
		Status:     "pending",
		Links:      []string{},
		TotalLinks: 0,
		CreatedAt:  time.Now(),
	}

	// Save task
	err = s.storage.SaveTask(task)
	if err != nil {
		return "", fmt.Errorf("failed to save task: %v", err)
	}

	// Start crawling in a goroutine
	go s.doCrawl(taskID, startURL, concurrentRequests, maxDepth, timeout)

	return taskID, nil
}

// doCrawl performs the actual crawling
func (s *CrawlerService) doCrawl(taskID, startURL string, concurrentRequests, maxDepth int, timeout time.Duration) {
	// Update task status
	err := s.storage.UpdateTaskStatus(taskID, "in_progress")
	if err != nil {
		fmt.Printf("Failed to update task status: %v\n", err)
		return
	}

	// Create crawler
	crawlerConfig := crawler.CrawlerConfig{
		StartURL:           startURL,
		ConcurrentRequests: concurrentRequests,
		MaxDepth:           maxDepth,
		Timeout:            timeout,
	}

	c, err := crawler.NewCrawler(crawlerConfig)
	if err != nil {
		s.storage.SetTaskFailed(taskID, err)
		fmt.Printf("Failed to create crawler: %v\n", err)
		return
	}

	// Start crawling
	startTime := time.Now()
	urlChan, err := c.Start()
	if err != nil {
		s.storage.SetTaskFailed(taskID, err)
		fmt.Printf("Failed to start crawler: %v\n", err)
		return
	}

	// Collect URLs from channel
	var urls []string
	for url := range urlChan {
		urls = append(urls, url)

		// Batch update to storage
		if len(urls)%100 == 0 {
			err := s.storage.AddTaskLinks(taskID, urls)
			if err != nil {
				fmt.Printf("Failed to add links to task: %v\n", err)
			}
			urls = []string{}
		}
	}

	// Add remaining URLs
	if len(urls) > 0 {
		err := s.storage.AddTaskLinks(taskID, urls)
		if err != nil {
			fmt.Printf("Failed to add links to task: %v\n", err)
		}
	}

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	// Mark task as completed
	err = s.storage.SetTaskCompleted(taskID, elapsedTime)
	if err != nil {
		fmt.Printf("Failed to mark task as completed: %v\n", err)
	}
}

// GetTask retrieves a task by its ID
func (s *CrawlerService) GetTask(taskID string) (*storage.Task, error) {
	return s.storage.GetTask(taskID)
}
