package main

import (
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData  // Changed from map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

// addPageVisit safely adds a page visit to the map
// Returns true if this is the first visit to this page
func (cfg *config) addPageVisit(normalizedURL string, pageData PageData) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Check if page already exists
	if _, exists := cfg.pages[normalizedURL]; exists {
		// Already visited - don't update, just return false
		return false
	}

	// First visit - add to map
	cfg.pages[normalizedURL] = pageData
	return true
}
