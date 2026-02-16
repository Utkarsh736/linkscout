package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// writeCSVReport writes the crawl results to a CSV file
func writeCSVReport(pages map[string]PageData, filename string) error {
	// Create the CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("couldn't create file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("couldn't write header: %w", err)
	}

	// Write data rows
	for _, pageData := range pages {
		// Join slices with semicolons
		outgoingLinks := strings.Join(pageData.OutgoingLinks, ";")
		imageURLs := strings.Join(pageData.ImageURLs, ";")

		// Create row
		row := []string{
			pageData.URL,
			pageData.H1,
			pageData.FirstParagraph,
			outgoingLinks,
			imageURLs,
		}

		// Write row to CSV
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("couldn't write row: %w", err)
		}
	}

	// Check for any errors during writing
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error writing CSV: %w", err)
	}

	return nil
}

