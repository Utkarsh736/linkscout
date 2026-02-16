package main

import "testing"

// ============================================
// Tests for getH1FromHTML
// ============================================

func TestGetH1FromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTMLMultiple(t *testing.T) {
	// Should return the FIRST h1
	inputBody := "<html><body><h1>First Title</h1><h1>Second Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "First Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTMLNoH1(t *testing.T) {
	inputBody := "<html><body><h2>Not an H1</h2></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTMLNestedTags(t *testing.T) {
	// H1 with nested span
	inputBody := "<html><body><h1>Title with <span>nested</span> tags</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Title with nested tags"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

// ============================================
// Tests for getFirstParagraphFromHTML
// ============================================

func TestGetFirstParagraphFromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><p>First paragraph.</p></body></html>"
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
<p>Outside paragraph.</p>
<main>
<p>Main paragraph.</p>
</main>
</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLNoMain(t *testing.T) {
	// Fallback to first <p> when no <main>
	inputBody := `<html><body>
<p>First paragraph.</p>
<p>Second paragraph.</p>
</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLNoParagraph(t *testing.T) {
	inputBody := "<html><body><div>No paragraph here</div></body></html>"
	actual := getFirstParagraphFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLEmptyMain(t *testing.T) {
	// Main exists but empty, fallback to outside <p>
	inputBody := `<html><body>
<p>Outside paragraph.</p>
<main></main>
</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Outside paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

