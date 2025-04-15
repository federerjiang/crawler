package api

import (
	"github.com/crawler/config"
	"github.com/crawler/internal/api/handlers"
	"github.com/crawler/internal/api/middleware"
	"github.com/crawler/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the API router
func SetupRouter(crawlerService *service.CrawlerService, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Use logger middleware
	router.Use(middleware.Logger())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Crawl endpoints
		crawlHandler := handlers.NewCrawlHandler(crawlerService, cfg)
		v1.POST("/crawl", crawlHandler.StartCrawl)

		// Task endpoints
		taskHandler := handlers.NewTaskHandler(crawlerService)
		v1.GET("/tasks/:taskId", middleware.ValidateTaskID(), taskHandler.GetTask)
	}

	return router
}
