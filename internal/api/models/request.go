package models

// CrawlRequest represents the request body for starting a crawl task
type CrawlRequest struct {
	StartURL           string `json:"start_url" binding:"required,url"`
	ConcurrentRequests int    `json:"concurrent_requests,omitempty"`
	MaxDepth           int    `json:"max_depth,omitempty"`
	Timeout            int    `json:"timeout,omitempty"`
}
