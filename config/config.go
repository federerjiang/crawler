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
	RandomUserAgent           bool          `yaml:"random_user_agent"`
	UserAgents                []string      `yaml:"user_agents"`
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
			RandomUserAgent:           false,
			UserAgents: []string{
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
				"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0",
			},
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

	if randomUA := os.Getenv("CRAWLER_RANDOM_USER_AGENT"); randomUA != "" {
		config.Crawler.RandomUserAgent = randomUA == "true" || randomUA == "1"
	}

	return config, nil
}
