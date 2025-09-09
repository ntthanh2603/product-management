# User Service

NestJS-based User Service with gRPC communication for microservice architecture.

## Overview

This service manages user data and provides gRPC endpoints to communicate with other services in the microservice system. The service is built with NestJS framework using TypeScript and uses SQLite database for storage.

## Technologies Used

- **NestJS**: Backend framework
- **TypeScript**: Programming language
- **gRPC**: Inter-service communication
- **SQLite**: Database
- **TypeORM**: ORM for database operations
- **Protocol Buffers**: Data serialization

## Project Structure

```
user-service/
├── src/
│   ├── config/          # Configuration files
│   ├── database/        # Database module and configuration
│   ├── proto/           # Generated protobuf types
│   ├── user/            # User module
│   ├── product/         # Product module (for cross-service communication)
│   ├── app.module.ts    # Main application module
│   └── main.ts          # Application entry point
├── .env.example         # Environment variables template
├── package.json         # Dependencies and scripts
├── Taskfile.yml         # Task runner commands
└── tsconfig.json        # TypeScript configuration
```

## Installation and Setup

### 1. Prerequisites

- Node.js (v18 or higher)
- npm or yarn package manager
- Protocol Buffers compiler (protoc)

### 2. Clone and Install Dependencies

```bash
# Navigate to user service directory
cd user-service

# Install dependencies
npm install

# Copy environment file
cp .env.example .env
```

### 3. Environment Configuration

Create a `.env` file with the following variables:

```env
NODE_ENV=development
PORT=3001
GRPC_PORT=50051
GRPC_URL=localhost:50051
DATABASE_PATH=users.db

# Product Service gRPC connection
PRODUCT_SERVICE_URL=localhost:50052
```

### 4. Generate Protobuf Types

```bash
# Generate TypeScript types from proto files
npm run proto

# Or use task runner
task proto
```

## Running the Service

### Development Mode

```bash
# Using npm
npm run start:dev

# Or use task runner
task dev
```

### Production Mode

```bash
# Build application
npm run build

# Start production server
npm run start:prod

# Or use task runner
task build
task start
```

### Using Task Runner

This service has a comprehensive Taskfile.yml with many useful commands:

```bash
# Development
task dev                 # Start development server with watch mode
task dev:debug          # Start with debug mode

# Building
task build              # Build application
task build:watch        # Build with watch mode

# Testing
task test               # Run unit tests
task test:watch         # Run tests with watch mode
task test:cov           # Run tests with coverage
task test:e2e           # Run end-to-end tests

# Code Quality
task lint               # Run linter
task format             # Format code with prettier
task typecheck          # Run TypeScript type checking

# Database
task db:migrate         # Run database migrations
task db:seed            # Seed database with test data

# gRPC & Proto
task proto              # Generate protobuf types
task grpc:health        # Check gRPC service health

# Dependencies
task deps:install       # Install dependencies
task deps:update        # Update dependencies
task deps:audit         # Security audit

# Utilities
task clean              # Clean build artifacts
task logs               # View service logs
task info               # Show service information
```

## gRPC Services

This service exposes the following gRPC methods:

### User Service

- `CreateUser`: Create a new user
- `GetUser`: Get user by ID
- `UpdateUser`: Update user information
- `DeleteUser`: Delete user
- `ListUsers`: Get list of users

### Health Check

- Service health check endpoint available at `grpc://localhost:50051`

## Database

The service uses SQLite database with TypeORM:

- Database file: `users.db`
- Auto-create tables when service starts
- Supports migrations and seeding

## API Testing

To test gRPC endpoints, you can use:

```bash
# Test with grpcurl
grpcurl -plaintext localhost:50051 list

# Or use built-in health check
task grpc:health
```

## Development Workflow

1. **Setup**: Install dependencies and setup environment
2. **Generate Proto**: Run `task proto` to generate types from proto files
3. **Development**: Use `task dev` to start development server
4. **Testing**: Run `task test` to verify functionality
5. **Linting**: Use `task lint` to check code quality
6. **Building**: Run `task build` to build for production

## Troubleshooting

### Common Issues

1. **Proto generation fails**

   ```bash
   # Ensure protoc is installed and proto files exist
   task proto:check
   ```

2. **Database connection issues**

   ```bash
   # Check database file permissions
   task db:check
   ```

3. **gRPC port conflicts**
   ```bash
   # Check if port is available
   task ports:check
   ```

### Logs

```bash
# View real-time logs
task logs

# View service information
task info
```

## Integration with Other Services

This service communicates with:

- **Product Service**: gRPC client connection
- **API Gateway**: Through gRPC server endpoints

Ensure all services are running before testing integration.
