package middleware

import (
	"net/http"

	"github.com/crawler/internal/api/models"
	"github.com/gin-gonic/gin"
)

// ValidateTaskID validates that the taskId is not empty
func ValidateTaskID() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("taskId")
		if taskID == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Error: &models.Error{
					Reason: "task ID is required",
					Field:  "taskId",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ValidateCrawlRequest validates the crawl request
func ValidateCrawlRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CrawlRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Error: &models.Error{
					Reason: err.Error(),
					Field:  "request",
				},
			})
			c.Abort()
			return
		}

		// Validate concurrent requests
		if request.ConcurrentRequests < 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Error: &models.Error{
					Reason: "concurrent_requests must be a positive integer",
					Field:  "concurrent_requests",
				},
			})
			c.Abort()
			return
		}

		// Validate max depth
		if request.MaxDepth < 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Error: &models.Error{
					Reason: "max_depth must be a positive integer",
					Field:  "max_depth",
				},
			})
			c.Abort()
			return
		}

		// Validate timeout
		if request.Timeout < 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Error: &models.Error{
					Reason: "timeout must be a positive integer",
					Field:  "timeout",
				},
			})
			c.Abort()
			return
		}

		// Set the validated request in the context
		c.Set("crawlRequest", request)
		c.Next()
	}
}
