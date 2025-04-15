package crawler

import (
	"sync"
)

// VisitedURLs represents a thread-safe structure for tracking visited URLs
type VisitedURLs struct {
	urls  map[string]bool
	mutex sync.RWMutex
}

// NewVisitedURLs creates a new instance of VisitedURLs
func NewVisitedURLs() *VisitedURLs {
	return &VisitedURLs{
		urls: make(map[string]bool),
	}
}

// Visit marks a URL as visited and returns true if the URL was not visited before
func (v *VisitedURLs) Visit(url string) bool {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if _, exists := v.urls[url]; exists {
		return false
	}

	v.urls[url] = true
	return true
}

// IsVisited checks if a URL has been visited
func (v *VisitedURLs) IsVisited(url string) bool {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	_, exists := v.urls[url]
	return exists
}

// GetVisitedURLs returns a slice of all visited URLs
func (v *VisitedURLs) GetVisitedURLs() []string {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	urls := make([]string, 0, len(v.urls))
	for url := range v.urls {
		urls = append(urls, url)
	}
	return urls
}

// Count returns the number of visited URLs
func (v *VisitedURLs) Count() int {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	return len(v.urls)
}
