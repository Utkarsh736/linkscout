package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

func main() {
	// Get command line arguments
	args := os.Args[1:]

	// Validate number of arguments
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	// Parse the base URL
	rawBaseURL := args[0]
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Error parsing base URL: %v\n", err)
		os.Exit(1)
	}

	// Print start message
	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	// Configure the crawler
	maxConcurrency := 5    // Number of concurrent requests
	maxPages := 100        // Maximum pages to crawl (safety limit)

	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	// Start the first crawl
	cfg.wg.Add(1)
	go func() {
		defer cfg.wg.Done()
		cfg.concurrencyControl <- struct{}{} // Acquire semaphore
		defer func() { <-cfg.concurrencyControl }() // Release semaphore
		
		cfg.crawlPage(rawBaseURL)
	}()

	// Wait for all goroutines to finish
	cfg.wg.Wait()

	// Print the results
	fmt.Println("\n=============================")
	fmt.Println("CRAWL COMPLETE")
	fmt.Println("=============================")
	fmt.Printf("Found %d unique pages:\n\n", len(cfg.pages))

	for url, count := range cfg.pages {
		fmt.Printf("%d - %s\n", count, url)
	}
}
