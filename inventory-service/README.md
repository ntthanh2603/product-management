# Inventory Service

A comprehensive inventory and order management microservice built with Python, gRPC, SQLAlchemy, and Kafka integration.

## 🎯 Features

### Inventory Management
- **Multi-location inventory tracking** across warehouses
- **Real-time stock availability** checking
- **Stock reservation system** with automatic expiration
- **Stock release** for cancelled orders
- **Concurrent stock operations** with proper locking

### Order Management
- **Complete order lifecycle** management
- **Multi-item orders** with validation
- **Order status tracking** (PENDING → CONFIRMED → PROCESSING → SHIPPED → DELIVERED)
- **Automatic stock validation** before order creation

### Event Streaming
- **Kafka integration** for real-time events
- **Order events**: Creation, status updates, cancellations
- **Inventory events**: Stock updates, reservations, releases
- **Event-driven architecture** for microservice communication

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                Inventory Service                        │
│                                                         │
│  ┌─────────────────┐    ┌─────────────────────────────┐ │
│  │  gRPC Server    │    │      Kafka Integration      │ │
│  │  Port: 50053    │    │                             │ │
│  │                 │    │  Producer: Send Events      │ │
│  │ • Inventory API │    │  Consumer: Process Events   │ │
│  │ • Order API     │    │                             │ │
│  └─────────────────┘    └─────────────────────────────┘ │
│           │                          │                  │
│           ▼                          ▼                  │
│  ┌─────────────────────────────────────────────────────┐ │
│  │              Business Logic                         │ │
│  │                                                     │ │
│  │ • Stock Management    • Order Processing           │ │
│  │ • Reservation Logic   • Status Tracking            │ │
│  │ • Validation Rules    • Event Handling             │ │
│  └─────────────────────────────────────────────────────┘ │
│           │                                             │
│           ▼                                             │
│  ┌─────────────────────────────────────────────────────┐ │
│  │                SQLite Database                      │ │
│  │                                                     │ │
│  │ • inventory_items     • orders                      │ │
│  │ • stock_reservations  • order_items                 │ │
│  └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Python 3.8+
- Task runner (optional but recommended)

### Installation

```bash
# Install dependencies
task install
# or
pip install -r requirements.txt

# Generate gRPC files
task proto

# Initialize database
task db-init
```

### Development

```bash
# Start the service
task dev

# Or run directly
python inventory_service.py
```

### Testing

```bash
# Run unit tests
task test

# Test gRPC interface
task test-grpc

# Run comprehensive system test
python ../test_inventory_system.py
```

## 📋 Available Tasks

### Development
```bash
task install      # Install dependencies
task proto        # Generate gRPC files
task dev          # Run in development mode
task health       # Check service health
task status       # Show service and database status
```

### Database Management
```bash
task db-init      # Initialize database
task db-reset     # Reset database
task db-shell     # Open SQLite shell
task seed-data    # Add sample data
task query-inventory  # Show current inventory
task query-orders     # Show current orders
```

### Kafka Operations
```bash
task kafka-topics           # List Kafka topics
task kafka-create-topics    # Create required topics
task kafka-consume-orders   # Monitor order events
task kafka-consume-inventory # Monitor inventory events
```

### Code Quality
```bash
task lint         # Lint code
task format       # Format code
task test         # Run tests
task clean        # Clean generated files
```

### Docker
```bash
task docker-build  # Build Docker image
task docker-run    # Run in container
task docker-stop   # Stop container
```

## 📊 Database Schema

### inventory_items
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| product_id | INTEGER | Product reference |
| quantity | INTEGER | Total quantity |
| reserved_quantity | INTEGER | Reserved quantity |
| location | STRING | Warehouse location |
| created_at | TIMESTAMP | Creation time |
| updated_at | TIMESTAMP | Last update time |

### orders
| Column | Type | Description |
|--------|------|-------------|
| id | STRING | UUID primary key |
| user_id | INTEGER | User reference |
| total_amount | DECIMAL | Order total |
| status | STRING | Order status |
| created_at | TIMESTAMP | Creation time |
| updated_at | TIMESTAMP | Last update time |

### order_items
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| order_id | STRING | Order reference |
| product_id | INTEGER | Product reference |
| quantity | INTEGER | Item quantity |
| price | DECIMAL | Item price |

### stock_reservations
| Column | Type | Description |
|--------|------|-------------|
| id | STRING | UUID primary key |
| product_id | INTEGER | Product reference |
| quantity | INTEGER | Reserved quantity |
| order_id | STRING | Order reference |
| is_active | BOOLEAN | Reservation status |
| created_at | TIMESTAMP | Creation time |
| expires_at | TIMESTAMP | Expiration time |

## 🔧 Configuration

### Environment Variables
```bash
DATABASE_URL=sqlite:///./inventory.db
GRPC_PORT=50053
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
```

### Docker Environment
```yaml
environment:
  - DATABASE_URL=sqlite:///./inventory.db
  - GRPC_PORT=50053
  - KAFKA_BOOTSTRAP_SERVERS=kafka:9092
```

## 📡 gRPC API Reference

### InventoryService

#### CreateInventoryItem
Create or update inventory for a product at a location.
```protobuf
rpc CreateInventoryItem(CreateInventoryItemRequest) returns (InventoryItemResponse);
```

#### GetInventoryItem
Retrieve inventory item by ID.
```protobuf
rpc GetInventoryItem(GetInventoryItemRequest) returns (InventoryItemResponse);
```

#### CheckStock
Check if sufficient stock is available.
```protobuf
rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
```

#### ReserveStock
Reserve stock for an order with automatic expiration.
```protobuf
rpc ReserveStock(ReserveStockRequest) returns (ReserveStockResponse);
```

#### ReleaseStock
Release previously reserved stock.
```protobuf
rpc ReleaseStock(ReleaseStockRequest) returns (ReleaseStockResponse);
```

### OrderService

#### CreateOrder
Create a new order with items and automatic stock validation.
```protobuf
rpc CreateOrder(CreateOrderRequest) returns (OrderResponse);
```

#### GetOrder
Retrieve order by ID.
```protobuf
rpc GetOrder(GetOrderRequest) returns (OrderResponse);
```

#### UpdateOrderStatus
Update order status and trigger events.
```protobuf
rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (OrderResponse);
```

#### ListOrders
List orders with pagination and filtering.
```protobuf
rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
```

## 🎪 Kafka Events

### Order Events (Topic: order-events)

#### ORDER_CREATED
```json
{
  "event_type": "ORDER_CREATED",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "id": "order-uuid",
    "user_id": 1,
    "total_amount": 199.99,
    "status": "PENDING",
    "items": [...]
  }
}
```

#### ORDER_CONFIRMED
```json
{
  "event_type": "ORDER_CONFIRMED",
  "timestamp": "2024-01-15T10:35:00Z",
  "data": {
    "id": "order-uuid",
    "status": "CONFIRMED"
  }
}
```

### Inventory Events (Topic: inventory-events)

#### STOCK_UPDATED
```json
{
  "event_type": "STOCK_UPDATED",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "product_id": 1,
    "quantity": 100,
    "location": "WAREHOUSE_A"
  }
}
```

#### STOCK_RESERVED
```json
{
  "event_type": "STOCK_RESERVED",
  "timestamp": "2024-01-15T10:31:00Z",
  "data": {
    "product_id": 1,
    "reserved_quantity": 10,
    "order_id": "order-uuid",
    "reservation_id": "res-uuid"
  }
}
```

## 🧪 Testing Examples

### Basic Inventory Operations
```bash
# Test gRPC interface
task test-grpc

# Check service health
task health

# View current inventory
task query-inventory
```

### Stock Reservation Workflow
```python
# Via gRPC client
import grpc
import inventory_pb2_grpc as pb2_grpc
import inventory_pb2 as pb2

channel = grpc.insecure_channel('localhost:50053')
client = pb2_grpc.InventoryServiceStub(channel)

# Check stock
response = client.CheckStock(pb2.CheckStockRequest(
    product_id=1,
    required_quantity=10
))

# Reserve stock
response = client.ReserveStock(pb2.ReserveStockRequest(
    product_id=1,
    quantity=10,
    order_id="order-123"
))
```

### Monitoring Events
```bash
# Monitor order events
task kafka-consume-orders

# Monitor inventory events  
task kafka-consume-inventory
```

## 🔍 Troubleshooting

### Common Issues

1. **gRPC Connection Failed**
   ```bash
   task health  # Check if service is running
   python inventory_service.py  # Start service manually
   ```

2. **Database Lock**
   ```bash
   task db-reset  # Reset database
   ```

3. **Kafka Connection Issues**
   ```bash
   docker ps | grep kafka  # Check Kafka status
   task kafka-topics       # List available topics
   ```

4. **Stock Reservation Problems**
   ```bash
   task query-inventory    # Check current stock levels
   # Reservations expire after 30 minutes automatically
   ```

### Logging
- Service logs to stdout with structured logging
- Kafka events are logged for debugging
- Database operations are logged (set echo=True in models.py)

### Performance Tips
- Use connection pooling for high throughput
- Consider partitioning Kafka topics by product_id
- Add database indexes for frequently queried fields
- Monitor reservation expiration cleanup

## 🔮 Future Enhancements

- [ ] **Distributed Locking**: Redis-based locking for high concurrency
- [ ] **Event Sourcing**: Store all state changes as events
- [ ] **CQRS**: Separate read and write models
- [ ] **Metrics**: Prometheus metrics integration
- [ ] **Saga Pattern**: Distributed transaction management
- [ ] **Multi-tenant**: Support for multiple organizations

## 📚 Related Documentation

- [Main Project README](../README.md)
- [Inventory System Overview](../INVENTORY_SYSTEM.md)
- [API Gateway Documentation](../api-gateway/README.md)
- [Docker Compose Setup](../docker-compose.yml)