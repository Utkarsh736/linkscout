# LinkScout 🕷️

A high-performance concurrent web crawler built in Go that analyzes websites and generates detailed CSV reports.

## Features

- ⚡ **Concurrent crawling** with configurable workers (goroutines)
- 🔒 **Thread-safe** with mutex protection and channels
- 📊 **CSV export** with page metadata (H1, paragraphs, links, images)
- 🎯 **Same-domain filtering** (won't crawl external sites)
- 🛑 **Configurable limits** (max pages, max concurrency)
- 🧪 **Fully tested** with comprehensive unit tests

## Installation

```bash
git clone https://github.com/Utkarsh736/linkscout.git
cd linkscout
go build -o crawler

