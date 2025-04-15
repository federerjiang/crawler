# Web Crawler Service

A high-performance web crawler service that allows you to crawl websites and extract links.

## Features

- Fast and concurrent web crawling
- RESTful API for starting crawl tasks and retrieving results
- Configurable crawling parameters (depth, timeout, concurrency)
- Domain-restricted crawling
- In-memory storage of crawl results
- Built-in Random User-Agent rotation with support for multiple browser types (Chrome, Firefox, Safari, Edge, Opera) and mobile devices

## Tech Stack

- Go 1.20+
- Gin Web Framework
- Colly (web scraping/crawling framework)

## Project Structure

```
.
├── cmd/
│   └── server/        # Entry point for the application
├── config/            # Configuration files
├── docs/              # Documentation
├── internal/
│   ├── api/           # API handlers, middleware, and models
│   ├── crawler/       # Crawler implementation
│   ├── service/       # Business logic
│   └── storage/       # Data storage implementation
├── pkg/
│   ├── logger/        # Logging functionality
│   └── utils/         # Utility functions
└── test/              # Unit and integration tests
```

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Git

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/crawler.git
cd crawler
```

2. Install dependencies:

```bash
go mod download
```

3. Build the project:

```bash
go build -o crawler ./cmd/server
```

4. Run the server:

```bash
./crawler
```

## Configuration

The application can be configured via the `config/config.yaml` file or environment variables:

```yaml
server:
  port: 8080
  timeout: 30s

crawler:
  default_concurrent_requests: 10
  default_timeout: 30s
  default_max_depth: 5
  random_user_agent: true
  user_agents:
    # Chrome
    - "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    # Firefox
    - "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0"
    # Safari
    - "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15"
    # Edge
    - "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59"
    # Mobile
    - "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1"
    # ...and many more
```

### Environment Variables

- `CONFIG_PATH`: Path to the configuration file
- `SERVER_PORT`: Server port
- `SERVER_TIMEOUT`: Server timeout
- `CRAWLER_CONCURRENT_REQUESTS`: Default number of concurrent requests
- `CRAWLER_TIMEOUT`: Default timeout for crawling requests
- `CRAWLER_MAX_DEPTH`: Default maximum depth for crawling
- `CRAWLER_RANDOM_USER_AGENT`: Enable/disable random user agent (true/false)

## API Reference

### Start a Crawl Task

```
POST /api/v1/crawl
```

Request body:

```json
{
  "start_url": "https://example.com",
  "concurrent_requests": 10,
  "max_depth": 3,
  "timeout": 30
}
```

Response:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "pending",
    "start_url": "https://example.com",
    "created_at": "2023-01-01T12:00:00Z"
  }
}
```

### Get Task Status

```
GET /api/v1/tasks/:taskId
```

Response:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "completed",
    "start_url": "https://example.com",
    "domain": "example.com",
    "links": ["https://example.com/page1", "https://example.com/page2"],
    "total_links": 2,
    "elapsed_time": "2.5s",
    "created_at": "2023-01-01T12:00:00Z",
    "completed_at": "2023-01-01T12:00:02Z"
  }
}
```

## Docker

Build the Docker image:

```bash
docker build -t crawler .
```

Run the container:

```bash
docker run -p 8080:8080 crawler
```

## License

MIT
