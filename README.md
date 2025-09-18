# Microservice System with NestJS & Python & Go

A complete microservice architecture featuring:

- **User Service**: NestJS with TypeScript and gRPC (following nest-grpc-base structure)
- **Product Service**: Python with gRPC and SQLAlchemy
- **Inventory Service**: Python with gRPC, SQLAlchemy, and Kafka integration
- **API Gateway**: Go Fiber REST API that communicates with all services via gRPC
- **Kafka**: Event streaming platform for real-time order and inventory events

## 🏗️ Architecture

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  User Service   │  │ Product Service │  │Inventory Service│
│   (NestJS)      │  │    (Python)     │  │    (Python)     │
│   Port: 50051   │  │   Port: 50052   │  │   Port: 50053   │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         ▲                     ▲                     ▲
         │ gRPC                │ gRPC                │ gRPC
         │                     │                     │
┌─────────────────────────────────────────────────────────────┐
│                 API Gateway (Go Fiber)                     │
│                    Port: 8000                              │
│               HTTP REST Endpoints                          │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
            ┌─────────────────────────────────┐
            │        Kafka Cluster            │
            │    (Zookeeper + Kafka)          │
            │   Ports: 2181, 9092            │
            │                                 │
            │ Topics:                         │
            │ • order-events                  │
            │ • inventory-events              │
            └─────────────────────────────────┘
```

## 📋 Services Overview

### User Service (NestJS)

- **Framework**: NestJS with TypeScript
- **Structure**: Following [nest-grpc-base](https://github.com/l1ttps/nest-grpc-base) repository pattern
- **Port**: 50051 (gRPC)
- **Database**: SQLite (`users.db`)
- **Features**:
  - User CRUD operations via gRPC
  - TypeORM for database management
  - Modular NestJS architecture
  - Proto-based type definitions
- **📖 Documentation**: [user-service/README.md](user-service/README.md)

### Product Service (Python)

- **Framework**: Python with gRPC and SQLAlchemy
- **Port**: 50052 (gRPC)
- **Database**: SQLite (`products.db`)
- **Features**:
  - Product CRUD operations
  - User validation (future: via gRPC to User Service)
  - SQLAlchemy ORM
  - Clean Python architecture
- **📖 Documentation**: [product-service/README.md](product-service/README.md)

### Inventory Service (Python)

- **Framework**: Python with gRPC, SQLAlchemy, and Kafka
- **Port**: 50053 (gRPC)
- **Database**: SQLite (`inventory.db`)
- **Features**:
  - Inventory item management with multi-location support
  - Stock reservation and release system
  - Order management with lifecycle tracking
  - Real-time Kafka event streaming
  - Automatic stock validation
  - Reservation expiration handling
- **📖 Documentation**: [INVENTORY_SYSTEM.md](INVENTORY_SYSTEM.md)

### API Gateway (Go Fiber)

- **Framework**: Go Fiber (high-performance web framework)
- **Port**: 8000 (HTTP REST)
- **Purpose**: Exposes gRPC services via HTTP REST API
- **Features**:
  - High-performance RESTful endpoints
  - Concurrent gRPC client connections
  - Built-in middleware (CORS, logging, recovery)
  - Fast JSON serialization
  - Error handling and validation
  - Inventory and order management endpoints
- **📖 Documentation**: [api-gateway/README.md](api-gateway/README.md)

### Kafka Event Streaming

- **Platform**: Apache Kafka with Zookeeper
- **Ports**: 2181 (Zookeeper), 9092 (Kafka)
- **Topics**:
  - `order-events`: Order lifecycle events
  - `inventory-events`: Stock and reservation events
- **Features**:
  - Real-time event streaming
  - Distributed messaging
  - Event-driven architecture
  - Microservice decoupling

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- Python 3.8+
- Go 1.21+
- npm

### Installation

1. **Install User Service (NestJS) dependencies:**

   ```bash
   cd user-service
   npm install
   ```

2. **Install Product Service (Python) dependencies:**

   ```bash
   cd product-service
   pip install -r requirements.txt
   ```

3. **Install API Gateway dependencies:**
   ```bash
   cd api-gateway
   go mod tidy
   ```

### Running Services

#### Option 1: Using Task Runner (Recommended) ⭐

**Install Task:**

```bash
# Windows (Chocolatey)
choco install go-task

# macOS (Homebrew)
brew install go-task/tap/go-task

# Linux (Snap)
snap install task --classic

# Or download from: https://taskfile.dev/installation/
```

**Run the project:**

```bash
# Install all dependencies
task install

# Start all services in development mode
task dev

# Run tests
task test

# Stop all services
task stop
```

#### Option 2: Manual Start

1. **Start User Service:**

   ```bash
   cd user-service
   npm run start:dev
   ```

2. **Start Product Service:**

   ```bash
   cd product-service
   python product-service.py
   ```

3. **Start API Gateway:**
   ```bash
   cd api-gateway
   go run main.go
   ```

#### Option 3: Using Scripts

**Windows:**

```bash
start-services.bat
```

**Linux/Mac:**

```bash
chmod +x start-services.sh
./start-services.sh
```

## 🧪 Testing

### Using Task Runner (Recommended) ⭐

```bash
# Run all tests
task test

# Test via gRPC only
task test-grpc

# Test via HTTP API Gateway only
task test-http

# Quick health check
task health
```

### Manual Testing

#### 1. gRPC Test Client (Python)

Test both services directly via gRPC:

```bash
python test_client.py
```

#### 2. HTTP REST API Testing (Go)

Test services via Go Fiber API Gateway:

```bash
go run test_client_go.go
```

### 3. HTTP REST API Testing (Manual)

If running the API Gateway, test via HTTP:

#### Create User:

```bash
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
  }'
```

#### Create Product:

```bash
curl -X POST http://localhost:8000/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "A test product",
    "price": 99.99,
    "user_id": 1
  }'
```

#### Get User's Products:

```bash
curl http://localhost:8000/api/users/1/products
```

## 📡 API Endpoints

### User Service Endpoints (via API Gateway)

| Method | Endpoint         | Description            |
| ------ | ---------------- | ---------------------- |
| POST   | `/api/users`     | Create a new user      |
| GET    | `/api/users`     | List users (paginated) |
| GET    | `/api/users/:id` | Get user by ID         |
| PUT    | `/api/users/:id` | Update user            |
| DELETE | `/api/users/:id` | Delete user            |

### Product Service Endpoints (via API Gateway)

| Method | Endpoint                  | Description               |
| ------ | ------------------------- | ------------------------- |
| POST   | `/api/products`           | Create a new product      |
| GET    | `/api/products`           | List products (paginated) |
| GET    | `/api/products/:id`       | Get product by ID         |
| PUT    | `/api/products/:id`       | Update product            |
| DELETE | `/api/products/:id`       | Delete product            |
| GET    | `/api/users/:id/products` | Get products by user      |

### Query Parameters

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)

## 🏛️ Project Structure

```
nest-microservice/
├── user-service/                 # NestJS User Service
│   ├── src/
│   │   ├── config/              # Configuration files
│   │   ├── database/            # Database entities & modules
│   │   ├── proto/               # Generated proto types
│   │   ├── product/             # Product gRPC client
│   │   ├── user/                # User module (controller, service)
│   │   ├── app.module.ts
│   │   └── main.ts
│   ├── package.json
│   ├── nest-cli.json
│   ├── tsconfig.json
│   ├── Taskfile.yml             # Service-specific tasks
│   └── README.md                # 📖 Service documentation
│
├── product-service/              # Python Product Service
│   ├── models.py                # SQLAlchemy models
│   ├── product-service.py       # gRPC service implementation
│   ├── product_pb2.py           # Generated proto
│   ├── product_pb2_grpc.py      # Generated gRPC stubs
│   ├── requirements.txt
│   ├── .env
│   ├── Taskfile.yml             # Service-specific tasks
│   └── README.md                # 📖 Service documentation
│
├── api-gateway/                 # Go Fiber API Gateway
│   ├── proto/                   # Generated Go proto files
│   ├── main.go                  # Fiber application
│   ├── go.mod                   # Go module file
│   ├── Taskfile.yml             # Service-specific tasks
│   └── README.md                # 📖 Service documentation
│
├── protos/                      # Protocol Buffer definitions
│   ├── user.proto
│   └── product.proto
│
├── test_client.py               # Python gRPC test client
├── test_client_go.go            # Go HTTP test client
├── start-services.bat/.sh       # Startup scripts
├── Taskfile.yml                 # Root task orchestrator
├── TASK_GUIDE.md               # Task usage guide
├── requirements.txt             # Gateway dependencies
└── README.md                   # Main project documentation
```

### 📚 Service Documentation

Each service has its own detailed README with setup and running instructions:

- **[User Service README](user-service/README.md)** - NestJS service setup, development workflow, and API reference
- **[Product Service README](product-service/README.md)** - Python service configuration, database management, and gRPC testing
- **[API Gateway README](api-gateway/README.md)** - Go Fiber gateway configuration, REST API endpoints, and performance tuning

## 🔧 Configuration

### User Service (.env)

```
DATABASE_URL=sqlite:./users.db
GRPC_PORT=50051
PRODUCT_SERVICE_URL=localhost:50052
```

### Product Service (.env)

```
DATABASE_URL=sqlite:///./products.db
GRPC_PORT=50052
USER_SERVICE_URL=localhost:50051
```

## 📊 Database Schema

### User (SQLite)

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Product (SQLite)

```sql
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔄 Communication Patterns

1. **Direct gRPC Communication:**

   - Client ↔ User Service (NestJS gRPC)
   - Client ↔ Product Service (Python gRPC)

2. **Service-to-Service Communication:**

   - User Service ↔ Product Service (via gRPC clients)

3. **HTTP REST via Gateway:**
   - Client ↔ API Gateway ↔ Services (gRPC)

## 🛠️ Development

### NestJS User Service

```bash
cd user-service
npm run start:dev    # Development with hot reload
npm run build        # Build for production
npm run proto        # Generate proto types
```

### Python Product Service

```bash
cd product-service
python product-service.py  # Start development server
```

### Generate Proto Files

```bash
# For NestJS (from user-service directory)
npm run proto

# For Python (from product-service directory)
python -m grpc_tools.protoc --proto_path=../protos --python_out=. --grpc_python_out=. ../protos/product.proto
```

## 🚨 Troubleshooting

### Common Issues

1. **Port Already in Use:**

   - Kill existing processes: `taskkill /f /im node.exe python.exe` (Windows)
   - Change ports in `.env` files

2. **Proto Import Errors:**

   - Regenerate proto files
   - Check Python path includes proto directories

3. **gRPC Connection Failed:**

   - Ensure services start in correct order (User → Product)
   - Check firewall settings

4. **Database Errors:**
   - Databases auto-create on first run
   - Delete `.db` files to reset

### Service Health Check

```bash
# Check API Gateway health
curl http://localhost:8000/health
```

## 🎯 Features Demonstrated

- ✅ NestJS microservice with gRPC (TypeScript)
- ✅ Python microservice with gRPC
- ✅ Go Fiber high-performance API Gateway
- ✅ Inter-service communication via gRPC
- ✅ Database integration (SQLite)
- ✅ HTTP REST API with concurrent gRPC clients
- ✅ Protocol Buffers for type-safe communication
- ✅ Error handling and middleware (CORS, logging, recovery)
- ✅ Structured logging across all services
- ✅ Modular architecture following best practices
- ✅ Multiple testing approaches (gRPC + HTTP)

## 📝 Next Steps

To extend this system:

1. Add authentication/authorization
2. Implement service discovery
3. Add monitoring and observability
4. Containerize with Docker
5. Add unit and integration tests
6. Implement circuit breakers
7. Add API rate limiting
8. Deploy to Kubernetes

This project demonstrates modern microservice patterns with different technologies working together seamlessly!

## ⚡ Task Runner Quick Reference

The project includes a comprehensive Task runner setup for easy management:

### 🚀 Essential Commands

```bash
task install    # Install all dependencies
task dev        # Start all services
task test       # Run all tests
task stop       # Stop all services
task health     # Check service health
```

### 🛠️ Development Commands

```bash
task dev-user      # Start only User Service
task dev-product   # Start only Product Service
task dev-gateway   # Start only API Gateway
task build         # Build all services
task lint          # Lint all code
```

### 🗑️ Cleanup Commands

```bash
task clean         # Clean databases and logs
task reset         # Full reset (clean + install + start)
task kill-processes # Force stop all processes
```

### 📚 Help

```bash
task               # Show all available tasks
task docs          # Show project documentation
```

For detailed Task usage, see **[TASK_GUIDE.md](TASK_GUIDE.md)**
