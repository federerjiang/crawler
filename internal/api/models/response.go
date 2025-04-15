package models

import "time"

// Response is the standard API response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error represents the error details in an error response
type Error struct {
	Reason string `json:"reason,omitempty"`
	Field  string `json:"field,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   *Error `json:"error,omitempty"`
}

// TaskResponse represents a crawl task response
type TaskResponse struct {
	TaskID      string    `json:"task_id"`
	Status      string    `json:"status"`
	StartURL    string    `json:"start_url"`
	Domain      string    `json:"domain,omitempty"`
	Links       []string  `json:"links,omitempty"`
	TotalLinks  int       `json:"total_links,omitempty"`
	ElapsedTime string    `json:"elapsed_time,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}
