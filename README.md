# LinkScout 🕷️

A high-performance concurrent web crawler built in Go that discovers internal links, extracts page metadata, and generates detailed CSV reports for SEO analysis and site mapping.

## Description

LinkScout was built to automate the process of mapping website structure and analyzing internal linking patterns. It recursively crawls web pages starting from a base URL, extracts key metadata (H1 tags, paragraphs, links, images), and exports everything to a structured CSV file — perfect for SEO audits, content analysis, or site documentation.

Built as part of the [Boot.dev Web Crawler course](https://www.boot.dev/courses/build-web-crawler-golang), this project showcases Go's concurrency primitives (goroutines, channels, mutexes, wait groups) and real-world software engineering practices.

## Features

- ⚡ **Concurrent crawling** with configurable worker pools (goroutines)
- 🔒 **Thread-safe** map access using `sync.Mutex`
- 🎯 **Domain filtering** - stays within the target domain (won't crawl external sites)
- 📊 **CSV export** with rich page metadata:
  - Page URL (normalized)
  - H1 tag content
  - First paragraph (prioritizes `<main>` content)
  - Outgoing links (semicolon-separated)
  - Image URLs (semicolon-separated)
- 🛑 **Configurable limits** - control max pages and concurrency
- 🧪 **Comprehensive test coverage** with table-driven tests
- 🚫 **Duplicate prevention** - tracks visited pages to avoid infinite loops
- 🌐 **User-Agent header** - identifies the crawler to avoid being blocked
- ✅ **Content-Type validation** - only crawls HTML pages, skips images/PDFs/etc.

## Prerequisites

- **Go 1.22 or higher**
- Internet connection (for crawling)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Utkarsh736/linkscout.git
   cd linkscout
   ```

2. Build the binary:
   ```bash
   go build -o crawler
   ```

3. _(Optional)_ Run tests:
   ```bash
   go test -v
   ```

## Usage

Run the crawler from the command line with three required arguments:

```bash
./crawler <URL> <maxConcurrency> <maxPages>
```

### Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `URL` | Base URL to start crawling (must include `http://` or `https://`) | `https://wagslane.dev` |
| `maxConcurrency` | Number of concurrent workers (1-20 recommended) | `5` |
| `maxPages` | Maximum pages to crawl before stopping | `50` |

### Examples

**Crawl a small blog:**
```bash
./crawler "https://wagslane.dev" 3 25
```

**Crawl a larger site with higher concurrency:**
```bash
./crawler "https://example.com" 10 100
```

**Single-threaded crawl for debugging:**
```bash
./crawler "https://example.com" 1 10
```

**Using `go run` instead of building:**
```bash
go run . "https://wagslane.dev" 5 50
```

## Output

After crawling completes, LinkScout generates **`report.csv`** in the current directory with the following structure:

| Column | Description | Example |
|--------|-------------|---------|
| `page_url` | Normalized URL | `wagslane.dev/posts/golang` |
| `h1` | H1 tag content | `"Learn Golang in 2026"` |
| `first_paragraph` | First paragraph text | `"Go is a statically typed..."` |
| `outgoing_link_urls` | Semicolon-separated links | `wagslane.dev/about;wagslane.dev/contact` |
| `image_urls` | Semicolon-separated images | `wagslane.dev/logo.png;wagslane.dev/banner.jpg` |

**Sample CSV:**
```csv
page_url,h1,first_paragraph,outgoing_link_urls,image_urls
wagslane.dev,Lane's Blog,Welcome to my blog,wagslane.dev/posts;wagslane.dev/about,wagslane.dev/logo.png
wagslane.dev/posts,All Posts,Here are my posts,wagslane.dev/posts/golang;wagslane.dev/posts/python,
wagslane.dev/about,About Me,I'm a software developer,,wagslane.dev/profile.jpg
```

Open `report.csv` in Excel, Google Sheets, or any CSV viewer for analysis.

## Architecture

```
linkscout/
├── main.go                  # CLI entry point, argument parsing
├── config.go                # Crawler configuration (mutex, channels, waitgroup)
├── crawl_page.go            # Recursive crawling logic with goroutines
├── fetch_html.go            # HTTP client with User-Agent headers
├── normalize_url.go         # URL normalization (remove schemes, trailing slashes)
├── get_html.go              # HTML parsing with goquery (H1, paragraphs)
├── get_urls.go              # Link and image extraction
├── page_data.go             # PageData struct and extraction logic
├── csv_report.go            # CSV export functionality
└── *_test.go                # Comprehensive unit tests
```

### Concurrency Design

LinkScout uses **Go's concurrency primitives** for safe, fast crawling:

- **Goroutines**: Each discovered link spawns a new goroutine
- **Buffered Channel (Semaphore)**: Limits max concurrent HTTP requests
- **Mutex**: Protects the shared `pages` map from race conditions
- **WaitGroup**: Ensures all goroutines complete before exiting

```go
// Simplified concurrency flow
cfg.wg.Add(1)
go func(url string) {
    defer cfg.wg.Done()
    cfg.concurrencyControl <- struct{}{}        // Acquire semaphore
    defer func() { <-cfg.concurrencyControl }() // Release semaphore
    
    cfg.mu.Lock()                               // Lock map
    cfg.pages[url] = pageData                   // Safe write
    cfg.mu.Unlock()                             // Unlock
    
    cfg.crawlPage(url)                          // Recurse
}(nextURL)
```

## Technologies

- **Go 1.22+** - Systems programming language with built-in concurrency
- **[goquery](https://github.com/PuerkitoBio/goquery)** - jQuery-like HTML parsing
- **net/http** - HTTP client with custom User-Agent headers
- **encoding/csv** - CSV file generation
- **sync.Mutex** - Thread-safe map access
- **Buffered Channels** - Semaphore pattern for concurrency control
- **sync.WaitGroup** - Goroutine synchronization

## Development

### Run Tests

```bash
# Run all tests
go test -v

# Run specific test
go test -run TestNormalizeURL

# Run with race detector
go run -race . "https://example.com" 3 10

# Check test coverage
go test -cover
```

### Key Test Files

- `normalize_url_test.go` - URL normalization edge cases
- `get_html_test.go` - HTML parsing (H1, paragraphs, main tags)
- `get_urls_test.go` - Link/image extraction and relative URL resolution
- `page_data_test.go` - PageData struct composition

### Debugging Tips

**View crawl progress in real-time:**
```bash
./crawler "https://example.com" 2 20 | tee crawl.log
```

**Kill stuck crawler:**
Press `Ctrl+C` to stop execution.

**Check for race conditions:**
```bash
go build -race -o crawler
./crawler "https://example.com" 5 25
```

## Project Background

This project was built as part of the **[Boot.dev Build a Web Crawler in Go course](https://www.boot.dev/courses/build-web-crawler-golang)**, which teaches:

- Test-driven development (TDD)
- Go concurrency patterns
- HTML parsing and HTTP clients
- Recursive algorithms
- Command-line tool development

## Future Improvements

- [ ] **JSON export** - Add JSON output option alongside CSV
- [ ] **Robots.txt compliance** - Respect site crawling rules
- [ ] **Rate limiting** - Add configurable delay between requests
- [ ] **Sitemap generation** - Export XML sitemap
- [ ] **Link graph visualization** - Generate network graph of page connections
- [ ] **External link tracking** - Count and report external links
- [ ] **Broken link detection** - Flag 404s and dead links
- [ ] **Progress bar** - Show real-time crawl progress
- [ ] **Docker support** - Containerize for easy deployment
- [ ] **Scheduled crawls** - Deploy with cron/scheduled tasks and email reports

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

**Utkarsh** - [GitHub](https://github.com/Utkarsh736)

---

⭐ **Star this repo if you find it helpful!**
