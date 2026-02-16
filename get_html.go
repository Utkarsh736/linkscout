package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	// Create reader from string
	reader := strings.NewReader(html)
	
	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ""
	}

	// Find first <h1> tag and get its text
	h1Text := doc.Find("h1").First().Text()
	
	return h1Text
}


func getFirstParagraphFromHTML(html string) string {
	reader := strings.NewReader(html)
	
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ""
	}

	// Try to find <main> tag first
	mainTag := doc.Find("main")
	
	var pText string
	
	// If <main> exists, look for <p> inside it
	if mainTag.Length() > 0 {
		pText = mainTag.Find("p").First().Text()
	}
	
	// If no <main> or no <p> in <main>, fallback to first <p> in document
	if pText == "" {
		pText = doc.Find("p").First().Text()
	}
	
	return pText
}


