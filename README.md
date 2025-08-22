# TheStartup HTTP Server

A custom HTTP/1.1 server implementation built from scratch in Go, featuring low-level TCP socket handling, HTTP protocol parsing, and advanced features like chunked transfer encoding with trailers.

## 🚀 Features

- **Custom HTTP/1.1 Server**: Full implementation without using Go's `net/http` package
- **Concurrent Request Handling**: Each connection handled in its own goroutine
- **HTTP Proxy**: Forward requests to external services with chunked transfer encoding
- **Video Streaming**: Serve MP4 video files with proper content-type headers
- **Custom Response Pages**: Beautifully crafted error and success pages
- **Graceful Shutdown**: Clean server shutdown with signal handling
- **Comprehensive Testing**: Unit tests for core components

## 📁 Project Structure

```
TheStartup/
├── cmd/                    # Application entry points
│   ├── httpserver/         # Main HTTP server
│   ├── tcplistener/        # TCP listener utility
│   └── udpsender/          # UDP sender utility
├── internal/               # Internal packages
│   ├── headers/            # HTTP headers parsing and management
│   ├── request/            # HTTP request parsing
│   ├── response/           # HTTP response writing
│   └── server/             # Core server implementation
├── assets/                 # Static assets
│   └── vim.mp4            # Example video file
└── go.mod                 # Go module definition
```

## 🛠️ Installation & Usage

### Prerequisites

- Go 1.24.2 or later

### Running the Server

```bash
# Clone the repository
git clone <repository-url>
cd TheStartup

# Install dependencies
go mod tidy

# Start the HTTP server
go run cmd/httpserver/main.go
```

The server will start on port `42069` and log startup confirmation.

### Graceful Shutdown

Press `Ctrl+C` or send a `SIGTERM` signal to gracefully stop the server.

## 🌐 API Endpoints

### Default Route
- **GET /** - Returns a success page with a friendly message
- **Response**: HTML page with "Your request was an absolute banger."

### Error Testing Routes
- **GET /yourproblem** - Returns a 400 Bad Request
- **GET /myproblem** - Returns a 500 Internal Server Error

### Media Serving
- **GET /video** - Serves the vim.mp4 video file
- **Content-Type**: `video/mp4`

### HTTP Proxy
- **GET /httpbin/*** - Proxies requests to `https://httpbin.org`
- **Features**: 
  - Chunked transfer encoding
  - SHA256 content verification trailers
  - Content-Length trailers

## 🔧 Technical Implementation

### Core Components

#### 1. HTTP Request Parser (`internal/request/`)
- Parses HTTP/1.1 request lines, headers, and body
- State-machine based parsing for streaming data
- Supports various HTTP methods and validates format
- Dynamic buffer resizing for large requests

#### 2. HTTP Response Writer (`internal/response/`)
- Structured response writing with state validation
- Support for chunked transfer encoding
- HTTP trailer implementation
- Content-type and content-length handling

#### 3. Headers Management (`internal/headers/`)
- Case-insensitive header storage and retrieval
- Header validation according to HTTP specifications
- Support for header overriding and removal
- Multi-value header handling

#### 4. Server Core (`internal/server/`)
- TCP listener with concurrent connection handling
- Graceful connection management
- Custom handler interface
- Thread-safe server shutdown

### Advanced Features

#### Chunked Transfer Encoding
The proxy handler implements proper chunked transfer encoding:

```go
// Streams data in 1024-byte chunks
const maxChunkSize = 1024

// Proper chunk termination
w.WriteChunkedBodyDone()

// Trailers with content verification
trailers.Override("X-Content-SHA256", sha256Hash)
trailers.Override("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
```

#### Content Integrity
- SHA256 hash calculation for proxied content
- Content-Length verification in trailers
- Error handling for corrupted transfers

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/headers/
go test ./internal/request/
```

### Example Test Commands

```bash
# Test the server manually
curl http://localhost:42069/
curl http://localhost:42069/yourproblem
curl http://localhost:42069/video
curl http://localhost:42069/httpbin/html
```

## 📊 Performance Characteristics

- **Concurrent Connections**: Unlimited (limited by system resources)
- **Memory Usage**: Dynamic buffer allocation with efficient reuse
- **Chunk Size**: 1024 bytes for optimal streaming performance
- **Connection Handling**: One goroutine per connection

## 🚧 Development

### Adding New Endpoints

1. Add route logic in `cmd/httpserver/main.go` handler function
2. Implement handler function following the signature:
   ```go
   func yourHandler(w *response.Writer, req *request.Request) {
       // Implementation
   }
   ```

### Custom Headers

```go
h := response.GetDefaultHeaders(contentLength)
h.Override("Custom-Header", "value")
h.Remove("Unwanted-Header")
w.WriteHeaders(h)
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## 📝 License

This project is part of the Boot.dev curriculum and is intended for educational purposes.

## 🙏 Acknowledgments

- Built as part of the [Boot.dev](https://boot.dev) Go course
- Implements HTTP/1.1 specification (RFC 7230-7235)
- Inspired by production HTTP servers like nginx and Apache

---

**Note**: This server is designed for educational purposes and may not be suitable for production use without additional security hardening and performance optimizations.
