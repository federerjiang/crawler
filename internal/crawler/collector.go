package crawler

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gocolly/colly/v2"
)

// CrawlerConfig represents the configuration for a crawler
type CrawlerConfig struct {
	StartURL           string
	ConcurrentRequests int
	MaxDepth           int
	Timeout            time.Duration
	UseRandomUserAgent bool
}

// Common User Agents
var userAgents = []string{
	// Chrome
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",

	// Firefox
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:92.0) Gecko/20100101 Firefox/92.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:93.0) Gecko/20100101 Firefox/93.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:90.0) Gecko/20100101 Firefox/90.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:91.0) Gecko/20100101 Firefox/91.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:91.0) Gecko/20100101 Firefox/91.0",

	// Safari
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15",

	// Edge
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36 Edg/92.0.902.78",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 Edg/93.0.961.38",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",

	// Opera
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 OPR/77.0.4054.254",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36 OPR/78.0.4093.147",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 OPR/77.0.4054.254",

	// Mobile Browsers
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 11; SM-G998U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.115 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 12; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.62 Mobile Safari/537.36",
}

// getRandomUserAgent returns a random user agent from the list
func getRandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
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
		// Set random user agent if enabled
		if c.config.UseRandomUserAgent {
			r.Headers.Set("User-Agent", getRandomUserAgent())
		}
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
