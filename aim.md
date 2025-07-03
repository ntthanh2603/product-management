# Học Microservice với gRPC và API Gateway trong NestJS

## Giới thiệu

Hướng dẫn này sẽ giúp bạn học cách xây dựng hệ thống microservice sử dụng gRPC để giao tiếp giữa các service và API Gateway REST để giao tiếp với client. Hướng dẫn này phù hợp cho những ai đã có kiến thức cơ bản về NestJS và NestJS Microservice với TCP.

## Kiến thức cần có

- NestJS cơ bản
- NestJS Microservice với TCP
- TypeScript
- Protocol Buffers (protobuf) cơ bản
- REST API concepts

## Mục lục

1. [Tổng quan Architecture](#1-tổng-quan-architecture)
2. [Cài đặt và cấu hình](#2-cài-đặt-và-cấu-hình)
3. [Tạo Proto Files](#3-tạo-proto-files)
4. [Xây dựng gRPC Microservices](#4-xây-dựng-grpc-microservices)
5. [Tạo API Gateway](#5-tạo-api-gateway)
6. [Giao tiếp giữa Services](#6-giao-tiếp-giữa-services)
7. [Error Handling](#7-error-handling)
8. [Testing](#8-testing)
9. [Deployment](#9-deployment)
10. [Best Practices](#10-best-practices)

## 1. Tổng quan Architecture

### Kiến trúc hệ thống

```
Client (REST API)
    ↓
API Gateway (REST)
    ↓ (gRPC)
┌─────────────────────────────────┐
│  User Service    │  Order Service │
│     (gRPC)       │     (gRPC)     │
└─────────────────────────────────┘
    ↓ (gRPC)           ↓ (gRPC)
┌─────────────────────────────────┐
│ Product Service  │ Payment Service│
│     (gRPC)       │     (gRPC)     │
└─────────────────────────────────┘
```

### Ưu điểm của gRPC so với TCP

- **Type Safety**: Strongly typed với Protocol Buffers
- **Performance**: Binary serialization, HTTP/2
- **Code Generation**: Tự động generate client/server code
- **Streaming**: Hỗ trợ bidirectional streaming
- **Cross-platform**: Hỗ trợ nhiều ngôn ngữ

## 2. Cài đặt và cấu hình

### Cài đặt dependencies

```bash
# Core packages
npm install @nestjs/microservices @grpc/grpc-js @grpc/proto-loader

# Development tools
npm install -D @types/node ts-proto
```

### Cấu trúc thư mục

```
project/
├── proto/
│   ├── user.proto
│   ├── order.proto
│   └── product.proto
├── apps/
│   ├── api-gateway/
│   ├── user-service/
│   ├── order-service/
│   └── product-service/
├── libs/
│   └── common/
└── package.json
```

## 3. Tạo Proto Files

### proto/user.proto

```protobuf
syntax = "proto3";

package user;

service UserService {
  rpc CreateUser (CreateUserRequest) returns (User);
  rpc GetUser (GetUserRequest) returns (User);
  rpc UpdateUser (UpdateUserRequest) returns (User);
  rpc DeleteUser (DeleteUserRequest) returns (DeleteUserResponse);
  rpc GetUsers (GetUsersRequest) returns (GetUsersResponse);
}

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
  string phone = 4;
  string created_at = 5;
  string updated_at = 6;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
  string phone = 3;
}

message GetUserRequest {
  int32 id = 1;
}

message UpdateUserRequest {
  int32 id = 1;
  string name = 2;
  string email = 3;
  string phone = 4;
}

message DeleteUserRequest {
  int32 id = 1;
}

message DeleteUserResponse {
  bool success = 1;
  string message = 2;
}

message GetUsersRequest {
  int32 page = 1;
  int32 limit = 2;
}

message GetUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}
```

### proto/order.proto

```protobuf
syntax = "proto3";

package order;

service OrderService {
  rpc CreateOrder (CreateOrderRequest) returns (Order);
  rpc GetOrder (GetOrderRequest) returns (Order);
  rpc GetOrdersByUser (GetOrdersByUserRequest) returns (GetOrdersByUserResponse);
  rpc UpdateOrderStatus (UpdateOrderStatusRequest) returns (Order);
}

message Order {
  int32 id = 1;
  int32 user_id = 2;
  repeated OrderItem items = 3;
  double total_amount = 4;
  string status = 5;
  string created_at = 6;
  string updated_at = 7;
}

message OrderItem {
  int32 product_id = 1;
  int32 quantity = 2;
  double price = 3;
}

message CreateOrderRequest {
  int32 user_id = 1;
  repeated OrderItem items = 2;
}

message GetOrderRequest {
  int32 id = 1;
}

message GetOrdersByUserRequest {
  int32 user_id = 1;
  int32 page = 2;
  int32 limit = 3;
}

message GetOrdersByUserResponse {
  repeated Order orders = 1;
  int32 total = 2;
}

message UpdateOrderStatusRequest {
  int32 id = 1;
  string status = 2;
}
```

## 4. Xây dựng gRPC Microservices

### User Service

#### apps/user-service/src/main.ts

```typescript
import { NestFactory } from "@nestjs/core";
import { MicroserviceOptions, Transport } from "@nestjs/microservices";
import { join } from "path";
import { UserModule } from "./user.module";

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    UserModule,
    {
      transport: Transport.GRPC,
      options: {
        package: "user",
        protoPath: join(__dirname, "../../../proto/user.proto"),
        url: "localhost:5001",
      },
    }
  );

  await app.listen();
  console.log("User Service is running on port 5001");
}
bootstrap();
```

#### apps/user-service/src/user.service.ts

```typescript
import { Injectable } from "@nestjs/common";
import { GrpcMethod } from "@nestjs/microservices";
import {
  CreateUserRequest,
  GetUserRequest,
  UpdateUserRequest,
  DeleteUserRequest,
  GetUsersRequest,
  User,
  DeleteUserResponse,
  GetUsersResponse,
} from "./interfaces/user.interface";

@Injectable()
export class UserService {
  private users: User[] = [];
  private idCounter = 1;

  @GrpcMethod("UserService", "CreateUser")
  createUser(request: CreateUserRequest): User {
    const user: User = {
      id: this.idCounter++,
      name: request.name,
      email: request.email,
      phone: request.phone,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    this.users.push(user);
    return user;
  }

  @GrpcMethod("UserService", "GetUser")
  getUser(request: GetUserRequest): User {
    const user = this.users.find((u) => u.id === request.id);
    if (!user) {
      throw new Error("User not found");
    }
    return user;
  }

  @GrpcMethod("UserService", "UpdateUser")
  updateUser(request: UpdateUserRequest): User {
    const userIndex = this.users.findIndex((u) => u.id === request.id);
    if (userIndex === -1) {
      throw new Error("User not found");
    }

    this.users[userIndex] = {
      ...this.users[userIndex],
      name: request.name,
      email: request.email,
      phone: request.phone,
      updatedAt: new Date().toISOString(),
    };

    return this.users[userIndex];
  }

  @GrpcMethod("UserService", "DeleteUser")
  deleteUser(request: DeleteUserRequest): DeleteUserResponse {
    const userIndex = this.users.findIndex((u) => u.id === request.id);
    if (userIndex === -1) {
      return {
        success: false,
        message: "User not found",
      };
    }

    this.users.splice(userIndex, 1);
    return {
      success: true,
      message: "User deleted successfully",
    };
  }

  @GrpcMethod("UserService", "GetUsers")
  getUsers(request: GetUsersRequest): GetUsersResponse {
    const page = request.page || 1;
    const limit = request.limit || 10;
    const startIndex = (page - 1) * limit;
    const endIndex = startIndex + limit;

    const users = this.users.slice(startIndex, endIndex);

    return {
      users,
      total: this.users.length,
    };
  }
}
```

#### apps/user-service/src/interfaces/user.interface.ts

```typescript
export interface User {
  id: number;
  name: string;
  email: string;
  phone: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  phone: string;
}

export interface GetUserRequest {
  id: number;
}

export interface UpdateUserRequest {
  id: number;
  name: string;
  email: string;
  phone: string;
}

export interface DeleteUserRequest {
  id: number;
}

export interface DeleteUserResponse {
  success: boolean;
  message: string;
}

export interface GetUsersRequest {
  page?: number;
  limit?: number;
}

export interface GetUsersResponse {
  users: User[];
  total: number;
}
```

### Order Service

#### apps/order-service/src/main.ts

```typescript
import { NestFactory } from "@nestjs/core";
import { MicroserviceOptions, Transport } from "@nestjs/microservices";
import { join } from "path";
import { OrderModule } from "./order.module";

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    OrderModule,
    {
      transport: Transport.GRPC,
      options: {
        package: "order",
        protoPath: join(__dirname, "../../../proto/order.proto"),
        url: "localhost:5002",
      },
    }
  );

  await app.listen();
  console.log("Order Service is running on port 5002");
}
bootstrap();
```

#### apps/order-service/src/order.service.ts

```typescript
import { Injectable } from "@nestjs/common";
import { GrpcMethod } from "@nestjs/microservices";
import {
  CreateOrderRequest,
  GetOrderRequest,
  GetOrdersByUserRequest,
  UpdateOrderStatusRequest,
  Order,
  GetOrdersByUserResponse,
} from "./interfaces/order.interface";

@Injectable()
export class OrderService {
  private orders: Order[] = [];
  private idCounter = 1;

  @GrpcMethod("OrderService", "CreateOrder")
  createOrder(request: CreateOrderRequest): Order {
    const totalAmount = request.items.reduce(
      (total, item) => total + item.price * item.quantity,
      0
    );

    const order: Order = {
      id: this.idCounter++,
      userId: request.userId,
      items: request.items,
      totalAmount,
      status: "pending",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    this.orders.push(order);
    return order;
  }

  @GrpcMethod("OrderService", "GetOrder")
  getOrder(request: GetOrderRequest): Order {
    const order = this.orders.find((o) => o.id === request.id);
    if (!order) {
      throw new Error("Order not found");
    }
    return order;
  }

  @GrpcMethod("OrderService", "GetOrdersByUser")
  getOrdersByUser(request: GetOrdersByUserRequest): GetOrdersByUserResponse {
    const userOrders = this.orders.filter((o) => o.userId === request.userId);
    const page = request.page || 1;
    const limit = request.limit || 10;
    const startIndex = (page - 1) * limit;
    const endIndex = startIndex + limit;

    const orders = userOrders.slice(startIndex, endIndex);

    return {
      orders,
      total: userOrders.length,
    };
  }

  @GrpcMethod("OrderService", "UpdateOrderStatus")
  updateOrderStatus(request: UpdateOrderStatusRequest): Order {
    const orderIndex = this.orders.findIndex((o) => o.id === request.id);
    if (orderIndex === -1) {
      throw new Error("Order not found");
    }

    this.orders[orderIndex] = {
      ...this.orders[orderIndex],
      status: request.status,
      updatedAt: new Date().toISOString(),
    };

    return this.orders[orderIndex];
  }
}
```

## 5. Tạo API Gateway

### apps/api-gateway/src/main.ts

```typescript
import { NestFactory } from "@nestjs/core";
import { ApiGatewayModule } from "./api-gateway.module";
import { ValidationPipe } from "@nestjs/common";

async function bootstrap() {
  const app = await NestFactory.create(ApiGatewayModule);

  app.useGlobalPipes(new ValidationPipe());
  app.enableCors();

  await app.listen(3000);
  console.log("API Gateway is running on port 3000");
}
bootstrap();
```

### apps/api-gateway/src/api-gateway.module.ts

```typescript
import { Module } from "@nestjs/common";
import { ClientsModule, Transport } from "@nestjs/microservices";
import { join } from "path";
import { UserController } from "./controllers/user.controller";
import { OrderController } from "./controllers/order.controller";

@Module({
  imports: [
    ClientsModule.register([
      {
        name: "USER_SERVICE",
        transport: Transport.GRPC,
        options: {
          package: "user",
          protoPath: join(__dirname, "../../../proto/user.proto"),
          url: "localhost:5001",
        },
      },
      {
        name: "ORDER_SERVICE",
        transport: Transport.GRPC,
        options: {
          package: "order",
          protoPath: join(__dirname, "../../../proto/order.proto"),
          url: "localhost:5002",
        },
      },
    ]),
  ],
  controllers: [UserController, OrderController],
})
export class ApiGatewayModule {}
```

### apps/api-gateway/src/controllers/user.controller.ts

```typescript
import {
  Controller,
  Get,
  Post,
  Put,
  Delete,
  Body,
  Param,
  Query,
  Inject,
  OnModuleInit,
} from "@nestjs/common";
import { ClientGrpc } from "@nestjs/microservices";
import { Observable } from "rxjs";
import {
  CreateUserDto,
  UpdateUserDto,
  GetUsersQueryDto,
} from "../dto/user.dto";

interface UserService {
  createUser(request: any): Observable<any>;
  getUser(request: any): Observable<any>;
  updateUser(request: any): Observable<any>;
  deleteUser(request: any): Observable<any>;
  getUsers(request: any): Observable<any>;
}

@Controller("users")
export class UserController implements OnModuleInit {
  private userService: UserService;

  constructor(@Inject("USER_SERVICE") private client: ClientGrpc) {}

  onModuleInit() {
    this.userService = this.client.getService<UserService>("UserService");
  }

  @Post()
  createUser(@Body() createUserDto: CreateUserDto) {
    return this.userService.createUser(createUserDto);
  }

  @Get(":id")
  getUser(@Param("id") id: string) {
    return this.userService.getUser({ id: parseInt(id) });
  }

  @Put(":id")
  updateUser(@Param("id") id: string, @Body() updateUserDto: UpdateUserDto) {
    return this.userService.updateUser({
      id: parseInt(id),
      ...updateUserDto,
    });
  }

  @Delete(":id")
  deleteUser(@Param("id") id: string) {
    return this.userService.deleteUser({ id: parseInt(id) });
  }

  @Get()
  getUsers(@Query() query: GetUsersQueryDto) {
    return this.userService.getUsers({
      page: query.page ? parseInt(query.page) : 1,
      limit: query.limit ? parseInt(query.limit) : 10,
    });
  }
}
```

### apps/api-gateway/src/dto/user.dto.ts

```typescript
import {
  IsString,
  IsEmail,
  IsNotEmpty,
  IsOptional,
  IsNumberString,
} from "class-validator";

export class CreateUserDto {
  @IsString()
  @IsNotEmpty()
  name: string;

  @IsEmail()
  @IsNotEmpty()
  email: string;

  @IsString()
  @IsNotEmpty()
  phone: string;
}

export class UpdateUserDto {
  @IsString()
  @IsOptional()
  name?: string;

  @IsEmail()
  @IsOptional()
  email?: string;

  @IsString()
  @IsOptional()
  phone?: string;
}

export class GetUsersQueryDto {
  @IsOptional()
  @IsNumberString()
  page?: string;

  @IsOptional()
  @IsNumberString()
  limit?: string;
}
```

## 6. Giao tiếp giữa Services

### Service-to-Service Communication

#### apps/order-service/src/order.service.ts (Updated)

```typescript
import { Injectable, Inject, OnModuleInit } from "@nestjs/common";
import { ClientGrpc, GrpcMethod } from "@nestjs/microservices";
import { Observable, firstValueFrom } from "rxjs";

interface UserService {
  getUser(request: any): Observable<any>;
}

@Injectable()
export class OrderService implements OnModuleInit {
  private userService: UserService;
  private orders: Order[] = [];
  private idCounter = 1;

  constructor(@Inject("USER_SERVICE") private userClient: ClientGrpc) {}

  onModuleInit() {
    this.userService = this.userClient.getService<UserService>("UserService");
  }

  @GrpcMethod("OrderService", "CreateOrder")
  async createOrder(request: CreateOrderRequest): Promise<Order> {
    // Verify user exists
    try {
      await firstValueFrom(this.userService.getUser({ id: request.userId }));
    } catch (error) {
      throw new Error("User not found");
    }

    const totalAmount = request.items.reduce(
      (total, item) => total + item.price * item.quantity,
      0
    );

    const order: Order = {
      id: this.idCounter++,
      userId: request.userId,
      items: request.items,
      totalAmount,
      status: "pending",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    this.orders.push(order);
    return order;
  }

  // ... other methods
}
```

#### apps/order-service/src/order.module.ts

```typescript
import { Module } from "@nestjs/common";
import { ClientsModule, Transport } from "@nestjs/microservices";
import { join } from "path";
import { OrderController } from "./order.controller";
import { OrderService } from "./order.service";

@Module({
  imports: [
    ClientsModule.register([
      {
        name: "USER_SERVICE",
        transport: Transport.GRPC,
        options: {
          package: "user",
          protoPath: join(__dirname, "../../../proto/user.proto"),
          url: "localhost:5001",
        },
      },
    ]),
  ],
  controllers: [OrderController],
  providers: [OrderService],
})
export class OrderModule {}
```

## 7. Error Handling

### Global Exception Filter

#### libs/common/src/filters/grpc-exception.filter.ts

```typescript
import { Catch, ArgumentsHost, HttpStatus } from "@nestjs/common";
import { BaseExceptionFilter } from "@nestjs/core";
import { RpcException } from "@nestjs/microservices";

@Catch(RpcException)
export class GrpcExceptionFilter extends BaseExceptionFilter {
  catch(exception: RpcException, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse();
    const request = ctx.getRequest();

    const error = exception.getError();

    let status = HttpStatus.INTERNAL_SERVER_ERROR;
    let message = "Internal server error";

    if (typeof error === "string") {
      message = error;
      if (error.includes("not found")) {
        status = HttpStatus.NOT_FOUND;
      } else if (error.includes("validation")) {
        status = HttpStatus.BAD_REQUEST;
      }
    } else if (typeof error === "object" && error !== null) {
      status = (error as any).status || HttpStatus.INTERNAL_SERVER_ERROR;
      message = (error as any).message || "Internal server error";
    }

    response.status(status).json({
      statusCode: status,
      message,
      timestamp: new Date().toISOString(),
      path: request.url,
    });
  }
}
```

### Sử dụng trong API Gateway

#### apps/api-gateway/src/main.ts (Updated)

```typescript
import { NestFactory } from "@nestjs/core";
import { ApiGatewayModule } from "./api-gateway.module";
import { ValidationPipe } from "@nestjs/common";
import { GrpcExceptionFilter } from "../../libs/common/src/filters/grpc-exception.filter";

async function bootstrap() {
  const app = await NestFactory.create(ApiGatewayModule);

  app.useGlobalPipes(new ValidationPipe());
  app.useGlobalFilters(new GrpcExceptionFilter());
  app.enableCors();

  await app.listen(3000);
  console.log("API Gateway is running on port 3000");
}
bootstrap();
```

## 8. Testing

### Unit Testing gRPC Service

#### apps/user-service/src/user.service.spec.ts

```typescript
import { Test, TestingModule } from "@nestjs/testing";
import { UserService } from "./user.service";
import { CreateUserRequest, GetUserRequest } from "./interfaces/user.interface";

describe("UserService", () => {
  let service: UserService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [UserService],
    }).compile();

    service = module.get<UserService>(UserService);
  });

  describe("createUser", () => {
    it("should create a user successfully", () => {
      const request: CreateUserRequest = {
        name: "John Doe",
        email: "john@example.com",
        phone: "123456789",
      };

      const result = service.createUser(request);

      expect(result).toBeDefined();
      expect(result.id).toBe(1);
      expect(result.name).toBe(request.name);
      expect(result.email).toBe(request.email);
      expect(result.phone).toBe(request.phone);
    });
  });

  describe("getUser", () => {
    it("should get a user successfully", () => {
      // Create a user first
      const createRequest: CreateUserRequest = {
        name: "John Doe",
        email: "john@example.com",
        phone: "123456789",
      };
      service.createUser(createRequest);

      const getRequest: GetUserRequest = { id: 1 };
      const result = service.getUser(getRequest);

      expect(result).toBeDefined();
      expect(result.id).toBe(1);
      expect(result.name).toBe(createRequest.name);
    });

    it("should throw error when user not found", () => {
      const getRequest: GetUserRequest = { id: 999 };

      expect(() => service.getUser(getRequest)).toThrow("User not found");
    });
  });
});
```

### Integration Testing

#### apps/api-gateway/src/controllers/user.controller.spec.ts

```typescript
import { Test, TestingModule } from "@nestjs/testing";
import { UserController } from "./user.controller";
import { ClientsModule, Transport } from "@nestjs/microservices";
import { join } from "path";

describe("UserController", () => {
  let controller: UserController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      imports: [
        ClientsModule.register([
          {
            name: "USER_SERVICE",
            transport: Transport.GRPC,
            options: {
              package: "user",
              protoPath: join(__dirname, "../../../../proto/user.proto"),
              url: "localhost:5001",
            },
          },
        ]),
      ],
      controllers: [UserController],
    }).compile();

    controller = module.get<UserController>(UserController);
  });

  it("should be defined", () => {
    expect(controller).toBeDefined();
  });
});
```

## 9. Deployment

### Docker Configuration

#### Dockerfile (User Service)

```dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

EXPOSE 5001

CMD ["node", "dist/apps/user-service/main"]
```

#### docker-compose.yml

```yaml
version: "3.8"

services:
  user-service:
    build:
      context: .
      dockerfile: apps/user-service/Dockerfile
    ports:
      - "5001:5001"
    environment:
      - NODE_ENV=production
    networks:
      - microservices-network

  order-service:
    build:
      context: .
      dockerfile: apps/order-service/Dockerfile
    ports:
      - "5002:5002"
    environment:
      - NODE_ENV=production
    depends_on:
      - user-service
    networks:
      - microservices-network

  api-gateway:
    build:
      context: .
      dockerfile: apps/api-gateway/Dockerfile
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - USER_SERVICE_URL=user-service:5001
      - ORDER_SERVICE_URL=order-service:5002
    depends_on:
      - user-service
      - order-service
    networks:
      - microservices-network

networks:
  microservices-network:
    driver: bridge
```

### Kubernetes Configuration

#### k8s/user-service-deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
        - name: user-service
          image: your-registry/user-service:latest
          ports:
            - containerPort: 5001
          env:
            - name: NODE_ENV
              value: "production"
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
    - port: 5001
      targetPort: 5001
  type: ClusterIP
```

## 10. Best Practices

### gRPC Best Practices

1. **Proto Design**

   - Sử dụng semantic versioning cho proto files
   - Định nghĩa rõ ràng message types
   - Sử dụng optional fields khi cần thiết

2. **Error Handling**

   - Implement proper error codes
   - Provide meaningful error messages
   - Use gRPC status codes correctly

3. **Performance**
   - Implement connection pooling
   - Use streaming for large data sets
   - Optimize message serialization

### Security

#### Implement Authentication

```typescript
// libs/common/src/guards/auth.guard.ts
import { Injectable, CanActivate, ExecutionContext } from "@nestjs/common";
import { Observable } from "rxjs";

@Injectable()
export class AuthGuard implements CanActivate {
  canActivate(
    context: ExecutionContext
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request = context.switchToHttp().getRequest();
    return this.validateToken(request.headers.authorization);
  }

  private validateToken(token: string): boolean {
    // Implement token validation logic
    return true;
  }
}
```

#### TLS Configuration

```typescript
// apps/user-service/src/main.ts (Updated)
import { NestFactory } from "@nestjs/core";
import { MicroserviceOptions, Transport } from "@nestjs/microservices";
import { join } from "path";
import { readFileSync } from "fs";

async function bootstrap() {
  const app = await NestFactory.createMicroservice<MicroserviceOptions>(
    UserModule,
    {
      transport: Transport.GRPC,
      options: {
        package: "user",
        protoPath: join(__dirname, "../../../proto/user.proto"),
        url: "localhost:5001",
        credentials: {
          key: readFileSync("./certs/server-key.pem"),
          cert: readFileSync("./certs/server-cert.pem"),
        },
      },
    }
  );

  await app.listen();
}
bootstrap();
```
