package handlers

import (
	"net/http"

	"github.com/crawler/internal/api/models"
	"github.com/crawler/internal/service"
	"github.com/gin-gonic/gin"
)

// TaskHandler handles the task API requests
type TaskHandler struct {
	crawlerService *service.CrawlerService
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(crawlerService *service.CrawlerService) *TaskHandler {
	return &TaskHandler{
		crawlerService: crawlerService,
	}
}

// GetTask handles the GET /api/v1/tasks/:taskId request
func (h *TaskHandler) GetTask(c *gin.Context) {
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
		return
	}

	// Get task
	task, err := h.crawlerService.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
			Error: &models.Error{
				Reason: "task not found",
			},
		})
		return
	}

	// Convert elapsed time to string
	elapsedTime := ""
	if task.ElapsedTime > 0 {
		elapsedTime = task.ElapsedTime.String()
	}

	// Return response
	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: models.TaskResponse{
			TaskID:      taskID,
			Status:      task.Status,
			StartURL:    task.StartURL,
			Domain:      task.Domain,
			Links:       task.Links,
			TotalLinks:  task.TotalLinks,
			ElapsedTime: elapsedTime,
			CreatedAt:   task.CreatedAt,
			CompletedAt: task.CompletedAt,
		},
	})
}
