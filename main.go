package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	// Get command line arguments
	args := os.Args[1:]

	// Validate number of arguments (now we need 3: URL, concurrency, maxPages)
	if len(args) < 3 {
		fmt.Println("not enough arguments provided")
		fmt.Println("usage: crawler <URL> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		fmt.Println("usage: crawler <URL> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}

	// Parse arguments
	rawBaseURL := args[0]
	
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("error parsing maxConcurrency: %v\n", err)
		os.Exit(1)
	}
	
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("error parsing maxPages: %v\n", err)
		os.Exit(1)
	}

	// Validate values
	if maxConcurrency < 1 {
		fmt.Println("maxConcurrency must be at least 1")
		os.Exit(1)
	}
	
	if maxPages < 1 {
		fmt.Println("maxPages must be at least 1")
		os.Exit(1)
	}

	// Parse the base URL
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing URL: %v\n", err)
		os.Exit(1)
	}

	// Print start message
	fmt.Printf("starting crawl of: %s\n", rawBaseURL)
	fmt.Printf("max concurrency: %d\n", maxConcurrency)
	fmt.Printf("max pages: %d\n", maxPages)
	fmt.Println()

	// Configure the crawler
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

	for pageURL, count := range cfg.pages {
		fmt.Printf("%d - %s\n", count, pageURL)
	}
}
