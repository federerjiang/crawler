package crawler

import (
	"errors"
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

// CrawlerConfig represents the configuration for a crawler
type CrawlerConfig struct {
	StartURL           string
	ConcurrentRequests int
	MaxDepth           int
	Timeout            time.Duration
}

// Crawler represents a crawler instance
type Crawler struct {
	config      CrawlerConfig
	collector   *colly.Collector
	visitedURLs *VisitedURLs
	domain      string
}

// NewCrawler creates a new crawler instance
func NewCrawler(config CrawlerConfig) (*Crawler, error) {
	// Extract domain from start URL
	domain, err := ExtractDomain(config.StartURL)
	if err != nil {
		return nil, fmt.Errorf("invalid start URL: %v", err)
	}

	// Create collector
	c := colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.MaxDepth(config.MaxDepth),
		colly.Async(true),
	)

	// Set timeout
	c.SetRequestTimeout(config.Timeout)

	// Set concurrent requests limit
	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: config.ConcurrentRequests,
		RandomDelay: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set limit rule: %v", err)
	}

	return &Crawler{
		config:      config,
		collector:   c,
		visitedURLs: NewVisitedURLs(),
		domain:      domain,
	}, nil
}

// Start starts the crawler and returns a channel for receiving URLs
func (c *Crawler) Start() (chan string, error) {
	// Normalize start URL
	normalizedURL, err := NormalizeURL(c.config.StartURL)
	if err != nil {
		return nil, fmt.Errorf("invalid start URL: %v", err)
	}

	urlChan := make(chan string, 100) // Buffer size to avoid blocking

	// Setup collector callbacks
	c.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}

		// Normalize the URL
		normalizedLink, err := NormalizeURL(link)
		if err != nil {
			return
		}

		// Check if the URL is from the same domain
		sameDomain, err := IsSameDomain(c.config.StartURL, normalizedLink)
		if err != nil || !sameDomain {
			return
		}

		// Mark URL as visited and add to channel if it's a new URL
		if c.visitedURLs.Visit(normalizedLink) {
			urlChan <- normalizedLink
			e.Request.Visit(normalizedLink)
		}
	})

	c.collector.OnRequest(func(r *colly.Request) {
		// Additional tracking if needed
	})

	c.collector.OnError(func(r *colly.Response, err error) {
		// Log errors
		fmt.Printf("Error visiting %s: %v\n", r.Request.URL, err)
	})

	// Visit the start URL
	if c.visitedURLs.Visit(normalizedURL) {
		urlChan <- normalizedURL
		err := c.collector.Visit(normalizedURL)
		if err != nil {
			close(urlChan)
			return nil, fmt.Errorf("failed to visit start URL: %v", err)
		}
	} else {
		close(urlChan)
		return nil, errors.New("start URL was already visited")
	}

	// Create a goroutine to close the channel when crawling is done
	go func() {
		c.collector.Wait()
		close(urlChan)
	}()

	return urlChan, nil
}

// Stop stops the crawler
func (c *Crawler) Stop() {
	// This method is a no-op for now as Colly doesn't provide a direct way to stop crawling
	// We could implement additional functionality here if needed
}

// GetVisitedURLs returns all visited URLs
func (c *Crawler) GetVisitedURLs() []string {
	return c.visitedURLs.GetVisitedURLs()
}

// GetVisitedCount returns the number of visited URLs
func (c *Crawler) GetVisitedCount() int {
	return c.visitedURLs.Count()
}
