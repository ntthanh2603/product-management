# Product Service

Python-based Product Service with gRPC communication for microservice architecture.

## Overview

This service manages product data and provides gRPC endpoints to communicate with other services in the microservice system. The service is built with Python using gRPC and uses SQLite database for storage.

## Technologies Used

- **Python 3.8+**: Programming language
- **gRPC**: Inter-service communication
- **SQLite**: Database
- **SQLAlchemy**: ORM for database operations
- **Protocol Buffers**: Data serialization
- **python-dotenv**: Environment variables management

## Project Structure

```
product-service/
├── models.py            # Database models (SQLAlchemy)
├── product_service.py   # Main gRPC service implementation
├── product_pb2.py       # Generated protobuf types
├── product_pb2_grpc.py  # Generated gRPC service stubs
├── requirements.txt     # Python dependencies
├── .env.example         # Environment variables template
└── Taskfile.yml         # Task runner commands
```

## Installation and Setup

### 1. Prerequisites

- Python 3.8 or higher
- pip package manager
- Protocol Buffers compiler (protoc)
- Virtual environment (recommended)

### 2. Clone and Install Dependencies

```bash
# Navigate to product service directory
cd product-service

# Create virtual environment
python -m venv venv

# Activate virtual environment
# Windows:
venv\Scripts\activate
# Linux/Mac:
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Copy environment file
cp .env.example .env
```

### 3. Environment Configuration

Create a `.env` file with the following variables:

```env
GRPC_PORT=50052
DATABASE_URL=sqlite:///products.db
LOG_LEVEL=INFO

# User Service gRPC connection (for validation)
USER_SERVICE_URL=localhost:50051
```

### 4. Generate Protobuf Types

```bash
# Generate Python protobuf types
python -m grpc_tools.protoc \
    --proto_path=../protos \
    --python_out=. \
    --grpc_python_out=. \
    ../protos/product.proto

# Or use task runner
task proto
```

## Running the Service

### Development Mode

```bash
# Using Python directly
python product_service.py

# Or use task runner
task dev
```

### Production Mode

```bash
# Install production dependencies
pip install -r requirements.txt

# Run service
task start

# Or with specific port
GRPC_PORT=50052 python product_service.py
```

### Using Task Runner

This service has a comprehensive Taskfile.yml with many useful commands:

```bash
# Development
task dev                 # Start development server
task dev:reload          # Start with auto-reload on file changes

# Virtual Environment
task venv:create         # Create virtual environment
task venv:activate       # Activate virtual environment
task venv:install        # Install dependencies in venv

# Dependencies
task deps:install        # Install dependencies
task deps:update         # Update dependencies
task deps:freeze         # Freeze current dependencies
task deps:check          # Check for security vulnerabilities

# Database
task db:init             # Initialize database
task db:migrate          # Run database migrations
task db:seed             # Seed database with test data
task db:reset            # Reset database

# Code Quality
task lint                # Run linting with flake8
task format              # Format code with black
task format:check        # Check formatting without changes
task type:check          # Run type checking with mypy

# Testing
task test                # Run tests with pytest
task test:unit           # Run unit tests only
task test:integration    # Run integration tests
task test:coverage       # Run tests with coverage report

# gRPC & Proto
task proto               # Generate protobuf types
task grpc:health         # Check gRPC service health
task grpc:test           # Test gRPC endpoints

# Building & Deployment
task build               # Build application
task package             # Package for deployment

# Utilities
task clean               # Clean temporary files
task logs                # View service logs
task info                # Show service information
task ports:check         # Check if ports are available
```

## gRPC Services

This service exposes the following gRPC methods:

### Product Service

- `CreateProduct`: Create a new product
- `GetProduct`: Get product by ID
- `UpdateProduct`: Update product information
- `DeleteProduct`: Delete product
- `ListProducts`: Get list of products with pagination
- `GetProductsByUser`: Get products by user ID

### Health Check

- Service health check endpoint available at `grpc://localhost:50052`

## Database Models

The service uses SQLAlchemy models:

```python
class Product:
    id: int (Primary Key)
    name: str
    description: str
    price: decimal
    user_id: int (Foreign Key to User Service)
    created_at: datetime
    updated_at: datetime
```

## API Examples

### Testing with grpcurl

```bash
# List available services
grpcurl -plaintext localhost:50052 list

# Create product
grpcurl -plaintext -d '{
  "name": "Test Product",
  "description": "Test Description",
  "price": 99.99,
  "user_id": 1
}' localhost:50052 product.ProductService/CreateProduct

# Get product
grpcurl -plaintext -d '{
  "product_id": 1
}' localhost:50052 product.ProductService/GetProduct

# List products with pagination
grpcurl -plaintext -d '{
  "page": 1,
  "limit": 10
}' localhost:50052 product.ProductService/ListProducts
```

### Testing with Task Runner

```bash
# Health check
task grpc:health

# Run integration tests
task grpc:test
```

## Development Workflow

1. **Setup**: Create virtual environment and install dependencies

   ```bash
   task venv:create
   task venv:activate
   task deps:install
   ```

2. **Generate Proto**: Run `task proto` to generate types from proto files

3. **Database Setup**: Initialize database

   ```bash
   task db:init
   task db:seed  # Optional: add test data
   ```

4. **Development**: Use `task dev` to start development server

5. **Testing**: Run tests to verify functionality

   ```bash
   task test
   task grpc:test
   ```

6. **Code Quality**: Check code quality
   ```bash
   task lint
   task format
   task type:check
   ```

## Database Management

```bash
# Initialize database tables
task db:init

# Reset database (careful in production!)
task db:reset

# Add sample data
task db:seed

# Check database status
task db:check
```

## Troubleshooting

### Common Issues

1. **Virtual environment issues**

   ```bash
   # Recreate venv
   task venv:create
   task venv:activate
   task deps:install
   ```

2. **Proto generation fails**

   ```bash
   # Check protoc installation
   protoc --version
   # Ensure proto files exist
   ls ../protos/
   ```

3. **Database connection issues**

   ```bash
   # Check database file
   task db:check
   # Reset database if needed
   task db:reset
   ```

4. **gRPC port conflicts**
   ```bash
   # Check if port is available
   task ports:check
   # Use different port
   GRPC_PORT=50053 task dev
   ```

### Dependencies Issues

```bash
# Check for security vulnerabilities
task deps:check

# Update outdated packages
task deps:update

# Reinstall clean environment
task clean
task venv:create
task deps:install
```

### Logs and Debugging

```bash
# View logs
task logs

# Run with debug mode
LOG_LEVEL=DEBUG task dev

# Service information
task info
```

## Integration with Other Services

This service communicates with:

- **User Service**: gRPC client connection to validate user_id
- **API Gateway**: Through gRPC server endpoints

### Cross-Service Communication

The service can call User Service to validate users:

```python
# TODO: Implement in CreateProduct method
# Validate user exists via gRPC call to user service
```

## Performance and Scaling

- **Connection Pooling**: SQLAlchemy connection pool
- **Thread Pool**: gRPC server with ThreadPoolExecutor (10 workers)
- **Database Indexing**: Indexes on user_id and created_at
- **Pagination**: Built-in pagination for ListProducts

## Deployment

```bash
# Build package
task build

# Package for deployment
task package

# Production dependencies only
pip install -r requirements.txt --no-dev
```

Ensure all required services (User Service, Database) are running before starting Product Service.
