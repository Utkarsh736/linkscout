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

	// Fetch the HTML from the current URL
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", rawCurrentURL, err)
		return
	}

	// Extract page data
	pageData := extractPageData(html, rawCurrentURL)

	// Check if this is the first visit to this page
	isFirst := cfg.addPageVisit(normalizedURL, pageData)
	if !isFirst {
		// Already visited - don't crawl again
		return
	}

	// Print progress
	fmt.Printf("Crawling: %s\n", rawCurrentURL)

	// Recursively crawl each URL found on the page (CONCURRENTLY!)
	for _, nextURL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go func(url string) {
			defer cfg.wg.Done()
			defer func() { <-cfg.concurrencyControl }() // Release semaphore

			cfg.concurrencyControl <- struct{}{} // Acquire semaphore
			cfg.crawlPage(url)
		}(nextURL)
	}
}
