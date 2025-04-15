package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/crawler/config"
	"github.com/crawler/internal/api"
	"github.com/crawler/internal/service"
	"github.com/crawler/internal/storage/inmemory"
	"github.com/crawler/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	configPath := "config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	// Make config path absolute if it's relative
	if !filepath.IsAbs(configPath) {
		dir, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current directory: %v", err)
			os.Exit(1)
		}
		configPath = filepath.Join(dir, configPath)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Create storage
	storage := inmemory.NewStorage()

	// Create services
	crawlerService := service.NewCrawlerService(storage)

	// Setup router
	router := api.SetupRouter(crawlerService, cfg)

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Starting server on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		logger.Error("Failed to start server: %v", err)
		os.Exit(1)
	}
}
