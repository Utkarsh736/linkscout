package main

import (
	"fmt"
	"net/url"
)

// crawlPage recursively crawls pages starting from rawCurrentURL
func (cfg *config) crawlPage(rawCurrentURL string) {
	// Check if we've reached max pages limit (thread-safe check)
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	// Parse current URL
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error parsing current URL: %v\n", err)
		return
	}

	// Check if current URL is on the same domain as base URL
	if currentURL.Host != cfg.baseURL.Host {
		// External link - don't crawl
		return
	}

	// Normalize the current URL
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing URL %s: %v\n", rawCurrentURL, err)
		return
	}

	// Check if this is the first visit to this page
	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		// Already visited - don't crawl again
		return
	}

	// Print progress (important for debugging!)
	fmt.Printf("Crawling: %s\n", rawCurrentURL)

	// Fetch the HTML from the current URL
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", rawCurrentURL, err)
		return
	}

	// Extract all URLs from the HTML
	urls, err := getURLsFromHTML(html, currentURL)
	if err != nil {
		fmt.Printf("Error extracting URLs from %s: %v\n", rawCurrentURL, err)
		return
	}

	// Recursively crawl each URL found on the page (CONCURRENTLY!)
	for _, nextURL := range urls {
		cfg.wg.Add(1)
		go func(url string) {
			defer cfg.wg.Done()
			defer func() { <-cfg.concurrencyControl }() // Release semaphore

			cfg.concurrencyControl <- struct{}{} // Acquire semaphore
			cfg.crawlPage(url)
		}(nextURL)
	}
}
