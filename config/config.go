package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Crawler CrawlerConfig `yaml:"crawler"`
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// CrawlerConfig holds the crawler configuration
type CrawlerConfig struct {
	DefaultConcurrentRequests int           `yaml:"default_concurrent_requests"`
	DefaultTimeout            time.Duration `yaml:"default_timeout"`
	DefaultMaxDepth           int           `yaml:"default_max_depth"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(configPath string) (*Config, error) {
	// Default configuration
	config := &Config{
		Server: ServerConfig{
			Port:    8080,
			Timeout: 30 * time.Second,
		},
		Crawler: CrawlerConfig{
			DefaultConcurrentRequests: 10,
			DefaultTimeout:            30 * time.Second,
			DefaultMaxDepth:           5,
		},
	}

	// Read configuration file
	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	// Parse YAML
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return config, err
	}

	// Override with environment variables if they exist
	if port := os.Getenv("SERVER_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil && p > 0 {
			config.Server.Port = p
		}
	}

	if timeout := os.Getenv("SERVER_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			config.Server.Timeout = t
		}
	}

	if concurrentRequests := os.Getenv("CRAWLER_CONCURRENT_REQUESTS"); concurrentRequests != "" {
		var c int
		if _, err := fmt.Sscanf(concurrentRequests, "%d", &c); err == nil && c > 0 {
			config.Crawler.DefaultConcurrentRequests = c
		}
	}

	if timeout := os.Getenv("CRAWLER_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			config.Crawler.DefaultTimeout = t
		}
	}

	if maxDepth := os.Getenv("CRAWLER_MAX_DEPTH"); maxDepth != "" {
		var d int
		if _, err := fmt.Sscanf(maxDepth, "%d", &d); err == nil && d > 0 {
			config.Crawler.DefaultMaxDepth = d
		}
	}

	return config, nil
}
