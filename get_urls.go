package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	// Parse HTML
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse HTML: %w", err)
	}

	var urls []string

	// Find all <a> tags with href attribute
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// Parse the href (could be relative or absolute)
		parsedHref, err := url.Parse(href)
		if err != nil {
			return // Skip invalid URLs
		}

		// Resolve relative URLs to absolute
		absoluteURL := baseURL.ResolveReference(parsedHref)
		
		urls = append(urls, absoluteURL.String())
	})

	return urls, nil
}


func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse HTML: %w", err)
	}

	var imageURLs []string

	// Find all <img> tags with src attribute
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		// Parse the src (could be relative or absolute)
		parsedSrc, err := url.Parse(src)
		if err != nil {
			return // Skip invalid URLs
		}

		// Resolve relative URLs to absolute
		absoluteURL := baseURL.ResolveReference(parsedSrc)
		
		imageURLs = append(imageURLs, absoluteURL.String())
	})

	return imageURLs, nil
}


