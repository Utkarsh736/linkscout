package main

import "net/url"

// PageData represents structured data extracted from a web page
type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

// extractPageData extracts and structures all relevant data from an HTML page
func extractPageData(html, pageURL string) PageData {
	// Parse the base URL for relative URL resolution
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		// If URL parsing fails, return minimal data
		return PageData{
			URL:           pageURL,
			OutgoingLinks: []string{},
			ImageURLs:     []string{},
		}
	}

	// Extract all data using our existing functions
	h1 := getH1FromHTML(html)
	firstParagraph := getFirstParagraphFromHTML(html)
	
	outgoingLinks, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		outgoingLinks = []string{}
	}
	
	imageURLs, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		imageURLs = []string{}
	}

	// Return structured data
	return PageData{
		URL:            pageURL,
		H1:             h1,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
