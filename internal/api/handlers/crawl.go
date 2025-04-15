package handlers

import (
	"net/http"
	"time"

	"github.com/crawler/config"
	"github.com/crawler/internal/api/models"
	"github.com/crawler/internal/service"
	"github.com/gin-gonic/gin"
)

// CrawlHandler handles the crawl API requests
type CrawlHandler struct {
	crawlerService *service.CrawlerService
	config         *config.Config
}

// NewCrawlHandler creates a new CrawlHandler
func NewCrawlHandler(crawlerService *service.CrawlerService, config *config.Config) *CrawlHandler {
	return &CrawlHandler{
		crawlerService: crawlerService,
		config:         config,
	}
}

// StartCrawl handles the POST /api/v1/crawl request
func (h *CrawlHandler) StartCrawl(c *gin.Context) {
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
		return
	}

	// Use default values if not provided
	concurrentRequests := request.ConcurrentRequests
	if concurrentRequests <= 0 {
		concurrentRequests = h.config.Crawler.DefaultConcurrentRequests
	}

	maxDepth := request.MaxDepth
	if maxDepth <= 0 {
		maxDepth = h.config.Crawler.DefaultMaxDepth
	}

	timeout := time.Duration(request.Timeout) * time.Second
	if timeout <= 0 {
		timeout = h.config.Crawler.DefaultTimeout
	}

	// 始终启用随机用户代理
	useRandomUserAgent := true

	// Start crawling
	taskID, err := h.crawlerService.StartCrawl(request.StartURL, concurrentRequests, maxDepth, timeout, useRandomUserAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error: &models.Error{
				Reason: err.Error(),
			},
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: models.TaskResponse{
			TaskID:    taskID,
			Status:    "pending",
			StartURL:  request.StartURL,
			CreatedAt: time.Now(),
		},
	})
}
