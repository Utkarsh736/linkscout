package main

import (
	"fmt"
	"os"
)

func main() {
	// Get command line arguments (excluding program name)
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

	// Extract the base URL
	baseURL := args[0]

	// Print start message
	fmt.Printf("starting crawl of: %s\n", baseURL)
}
