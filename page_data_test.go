package main

import (
	"reflect"
	"testing"
)

func TestExtractPageDataBasic(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
		<h1>Test Title</h1>
		<p>This is the first paragraph.</p>
		<a href="/link1">Link 1</a>
		<img src="/image1.jpg" alt="Image 1">
	</body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://blog.boot.dev",
		H1:             "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
		ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataMultipleLinks(t *testing.T) {
	inputURL := "https://example.com"
	inputBody := `<html><body>
		<h1>Multiple Links</h1>
		<p>First paragraph text.</p>
		<a href="/page1">Page 1</a>
		<a href="https://other.com/page2">Page 2</a>
		<a href="/page3">Page 3</a>
		<img src="/img1.png">
		<img src="/img2.jpg">
	</body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://example.com",
		H1:             "Multiple Links",
		FirstParagraph: "First paragraph text.",
		OutgoingLinks: []string{
			"https://example.com/page1",
			"https://other.com/page2",
			"https://example.com/page3",
		},
		ImageURLs: []string{
			"https://example.com/img1.png",
			"https://example.com/img2.jpg",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataWithMain(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
		<h1>Page Title</h1>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
		<a href="/link">Link</a>
	</body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://blog.boot.dev",
		H1:             "Page Title",
		FirstParagraph: "Main paragraph.", // Should prioritize <main>
		OutgoingLinks:  []string{"https://blog.boot.dev/link"},
		ImageURLs:      []string{},
	}

	// Compare fields individually for better error messages
	if actual.URL != expected.URL {
		t.Errorf("URL: expected %q, got %q", expected.URL, actual.URL)
	}
	if actual.H1 != expected.H1 {
		t.Errorf("H1: expected %q, got %q", expected.H1, actual.H1)
	}
	if actual.FirstParagraph != expected.FirstParagraph {
		t.Errorf("FirstParagraph: expected %q, got %q", expected.FirstParagraph, actual.FirstParagraph)
	}
	if !reflect.DeepEqual(actual.OutgoingLinks, expected.OutgoingLinks) {
		t.Errorf("OutgoingLinks: expected %v, got %v", expected.OutgoingLinks, actual.OutgoingLinks)
	}
	if len(actual.ImageURLs) != len(expected.ImageURLs) {
		t.Errorf("ImageURLs length: expected %d, got %d", len(expected.ImageURLs), len(actual.ImageURLs))
	}
}

func TestExtractPageDataMissingElements(t *testing.T) {
	inputURL := "https://example.com"
	inputBody := `<html><body>
		<div>No h1, no paragraph, no links, no images</div>
	</body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://example.com",
		H1:             "",
		FirstParagraph: "",
		OutgoingLinks:  []string{},
		ImageURLs:      []string{},
	}

	if actual.URL != expected.URL {
		t.Errorf("URL: expected %q, got %q", expected.URL, actual.URL)
	}
	if actual.H1 != expected.H1 {
		t.Errorf("H1: expected %q, got %q", expected.H1, actual.H1)
	}
	if actual.FirstParagraph != expected.FirstParagraph {
		t.Errorf("FirstParagraph: expected %q, got %q", expected.FirstParagraph, actual.FirstParagraph)
	}
	if len(actual.OutgoingLinks) != 0 {
		t.Errorf("OutgoingLinks: expected empty, got %v", actual.OutgoingLinks)
	}
	if len(actual.ImageURLs) != 0 {
		t.Errorf("ImageURLs: expected empty, got %v", actual.ImageURLs)
	}
}

func TestExtractPageDataRelativeURLs(t *testing.T) {
	inputURL := "https://blog.boot.dev/path/page"
	inputBody := `<html><body>
		<h1>Relative URLs</h1>
		<p>Testing relative paths.</p>
		<a href="../other">Parent Path</a>
		<a href="./same">Same Level</a>
		<img src="../images/logo.png">
	</body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://blog.boot.dev/path/page",
		H1:             "Relative URLs",
		FirstParagraph: "Testing relative paths.",
		OutgoingLinks: []string{
			"https://blog.boot.dev/other",
			"https://blog.boot.dev/path/same",
		},
		ImageURLs: []string{
			"https://blog.boot.dev/images/logo.png",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

