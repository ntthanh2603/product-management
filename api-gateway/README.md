# API Gateway

Go Fiber-based API Gateway for microservice architecture with gRPC backend communication.

## Overview

This API Gateway serves as the entry point for the entire microservice system, providing REST API endpoints and routing requests to appropriate gRPC services. The service is built with Go Fiber framework to achieve high performance and low latency.

## Technologies Used

- **Go 1.21+**: Programming language
- **Fiber v2**: High-performance web framework
- **gRPC**: Backend service communication
- **Protocol Buffers**: Data serialization
- **CORS**: Cross-Origin Resource Sharing support
- **Logging & Recovery**: Built-in middleware

## Project Structure

```
api-gateway/
├── main.go              # Main application file với routes và handlers
├── proto/               # Generated protobuf types và gRPC clients
├── go.mod               # Go module dependencies
├── go.sum               # Dependencies checksums
├── Taskfile.yml         # Task runner commands
└── README.md            # Documentation
```

## Architecture

```
HTTP REST Client → API Gateway (Fiber) → gRPC Services
                       ↓
               ┌─────────────────┐
               │   API Gateway   │
               │   (Port 8000)   │
               └─────────────────┘
                       ↓
          ┌─────────────────┬─────────────────┐
          ↓                 ↓                 ↓
   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐
   │ User Service│   │Product Svc  │   │Other Services│
   │(Port 50051) │   │(Port 50052) │   │             │
   └─────────────┘   └─────────────┘   └─────────────┘
```

## Installation and Setup

### 1. Prerequisites

- Go 1.21 or higher
- Protocol Buffers compiler (protoc)
- Go gRPC tools

### 2. Install Dependencies

```bash
# Navigate to api-gateway directory
cd api-gateway

# Download dependencies
go mod download

# Install protoc-gen-go tools if needed
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 3. Environment Configuration

The API Gateway uses the following default configurations:

```go
// Server Configuration
Port: 8000

// gRPC Service Endpoints
User Service: localhost:50051
Product Service: localhost:50052

// Timeouts
Context Timeout: 5 seconds
```

### 4. Generate Protobuf Types

```bash
# Generate Go protobuf types
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       ../protos/*.proto

# Or use task runner
task proto
```

## Running the Service

### Development Mode

```bash
# Using go run
go run main.go

# Or use task runner
task dev
```

### Production Mode

```bash
# Build binary
go build -o api-gateway

# Run binary
./api-gateway

# Or use task runner
task build
task start
```

### Using Task Runner

This service has a comprehensive Taskfile.yml with many useful commands:

```bash
# Development
task dev                 # Start development server with auto-reload
task dev:watch           # Start with file watching
task dev:debug           # Start with debug mode

# Building
task build               # Build binary
task build:linux         # Build for Linux
task build:windows       # Build for Windows
task build:mac           # Build for macOS
task build:all           # Build for all platforms

# Dependencies
task deps:download       # Download dependencies
task deps:update         # Update dependencies
task deps:tidy           # Clean up dependencies
task deps:vendor         # Vendor dependencies

# Code Quality
task lint                # Run linter (golangci-lint)
task format              # Format code with gofmt
task vet                 # Run go vet
task staticcheck         # Run static analysis

# Testing
task test                # Run unit tests
task test:race           # Run tests with race detection
task test:coverage       # Run tests with coverage
task test:bench          # Run benchmark tests

# Security
task security            # Run security checks
task vuln:check          # Check for vulnerabilities

# gRPC & Proto
task proto               # Generate protobuf types
task grpc:health         # Check gRPC services health

# Utilities
task clean               # Clean build artifacts
task info                # Show service information
task ports:check         # Check if ports are available
task services:check      # Check if backend services are running
```

## API Endpoints

### REST API Routes

#### Users

- `POST /api/users` - Create new user
- `GET /api/users` - List users (with pagination)
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user
- `GET /api/users/:id/products` - Get user's products

#### Products

- `POST /api/products` - Create new product
- `GET /api/products` - List products (with pagination)
- `GET /api/products/:id` - Get product by ID

#### Health Check

- `GET /health` - Service health status

### Request/Response Examples

#### Create User

```bash
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30
}

# Response
{
  "success": true,
  "message": "User created successfully",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### List Products with Pagination

```bash
GET /api/products?page=1&limit=10

# Response
{
  "products": [...],
  "total": 25,
  "page": 1,
  "limit": 10
}
```

#### Health Check

```bash
GET /health

# Response
{
  "status": "healthy",
  "time": "2024-01-01T00:00:00Z",
  "services": {
    "user-service": "localhost:50051",
    "product-service": "localhost:50052"
  }
}
```

## Testing API

### Using cURL

```bash
# Health check
curl http://localhost:8000/health

# Create user
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","age":25}'

# Get users
curl http://localhost:8000/api/users?page=1&limit=5

# Create product
curl -X POST http://localhost:8000/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Product","description":"Test Description","price":99.99,"user_id":1}'
```

### Using Task Runner

```bash
# Check if services are healthy
task grpc:health

# Test all endpoints
task test:integration
```

## Development Workflow

1. **Setup**: Install Go and dependencies

   ```bash
   go mod download
   task deps:download
   ```

2. **Generate Proto**: Create protobuf types

   ```bash
   task proto
   ```

3. **Development**: Start development server

   ```bash
   task dev
   ```

4. **Testing**: Run tests and check functionality

   ```bash
   task test
   task services:check
   ```

5. **Code Quality**: Ensure code quality

   ```bash
   task lint
   task format
   task vet
   ```

6. **Building**: Build for production
   ```bash
   task build
   # Or build for specific platform
   task build:linux
   ```

## Performance Features

- **High Performance**: Fiber framework built on Fasthttp
- **Connection Pooling**: Persistent gRPC connections
- **Concurrent Processing**: Goroutines for parallel processing
- **Middleware**: Built-in logging, CORS, recovery
- **Timeout Management**: Context-based request timeouts

## Monitoring and Logging

```go
// Logging format
${time} ${status} - ${method} ${path} ${latency}

// Example log output
2024-01-01T10:00:00Z 200 - GET /api/users 15ms
2024-01-01T10:00:01Z 201 - POST /api/products 25ms
```

## Error Handling

The API Gateway has comprehensive error handling:

- **Validation Errors**: 400 Bad Request
- **Not Found**: 404 Not Found
- **gRPC Errors**: Mapped to appropriate HTTP status codes
- **Internal Errors**: 500 Internal Server Error với detailed logging

## Security Features

- **CORS**: Configured for cross-origin requests
- **Input Validation**: Request body validation
- **Error Sanitization**: Hide internal details in production
- **Timeout Protection**: Prevent hanging requests

## Troubleshooting

### Common Issues

1. **gRPC connection fails**

   ```bash
   # Check if backend services are running
   task services:check

   # Check ports
   task ports:check
   ```

2. **Build failures**

   ```bash
   # Clean and rebuild
   task clean
   go mod tidy
   task build
   ```

3. **Proto generation issues**

   ```bash
   # Check protoc installation
   protoc --version

   # Regenerate protos
   task proto
   ```

### Debugging

```bash
# Run with debug mode
task dev:debug

# Check logs
task info

# Test individual services
curl http://localhost:8000/health
```

## Deployment

### Single Binary Deployment

```bash
# Build production binary
task build

# Copy binary to server
scp api-gateway user@server:/path/to/deploy/

# Run on server
./api-gateway
```

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-gateway

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-gateway .
EXPOSE 8000
CMD ["./api-gateway"]
```

### Multi-Platform Builds

```bash
# Build for multiple platforms
task build:all

# Output binaries:
# api-gateway-linux-amd64
# api-gateway-windows-amd64.exe
# api-gateway-darwin-amd64
```

## Integration with Backend Services

The API Gateway requires the following services to be running:

1. **User Service**: gRPC server trên port 50051
2. **Product Service**: gRPC server trên port 50052

### Service Discovery

Currently uses static configuration. Can be extended with:

- Consul service discovery
- etcd service registry
- Kubernetes service discovery

## Performance Tuning

### Recommended Settings

```go
// Production configuration
app := fiber.New(fiber.Config{
    Prefork: true,          // Enable prefork
    CaseSensitive: true,    // Enable case sensitivity
    StrictRouting: true,    // Enable strict routing
    ServerHeader: "Fiber", // Server header
    AppName: "API Gateway", // App name
})
```

### Monitoring Metrics

- Request/Response times
- Error rates
- gRPC connection health
- Memory usage
- Goroutine count

Ensure all backend services (User Service, Product Service) are running before starting the API Gateway.
