package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"api-gateway/proto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClients struct {
	UserClient    proto.UserServiceClient
	ProductClient proto.ProductServiceClient
}

var clients *GrpcClients

func main() {
	// Initialize gRPC clients
	var err error
	clients, err = initGrpcClients()
	if err != nil {
		log.Fatal("Failed to initialize gRPC clients:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: globalErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} ${status} - ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", healthCheck)

	// API routes
	api := app.Group("/api")

	// User routes
	userRoutes := api.Group("/users")
	userRoutes.Post("/", createUser)
	userRoutes.Get("/", listUsers)
	userRoutes.Get("/:id", getUser)
	userRoutes.Put("/:id", updateUser)
	userRoutes.Delete("/:id", deleteUser)
	userRoutes.Get("/:id/products", getUserProducts)

	// Product routes
	productRoutes := api.Group("/products")
	productRoutes.Post("/", createProduct)
	productRoutes.Get("/", listProducts)
	productRoutes.Get("/:id", getProduct)

	log.Println("🚀 API Gateway starting on port 8000")
	log.Println("📍 User endpoints: /api/users")
	log.Println("📍 Product endpoints: /api/products")
	log.Println("📍 Health check: /health")

	log.Fatal(app.Listen(":8000"))
}

func initGrpcClients() (*GrpcClients, error) {
	// Connect to User Service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Connect to Product Service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcClients{
		UserClient:    proto.NewUserServiceClient(userConn),
		ProductClient: proto.NewProductServiceClient(productConn),
	}, nil
}

func globalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
		"time":   time.Now(),
		"services": fiber.Map{
			"user-service":    "localhost:50051",
			"product-service": "localhost:50052",
		},
	})
}

// User endpoint handlers
func createUser(c *fiber.Ctx) error {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int32  `json:"age"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.UserClient.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": resp.Success,
		"message": resp.Message,
		"user":    resp.User,
	})
}

func getUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.UserClient.GetUser(ctx, &proto.GetUserRequest{
		UserId: int32(id),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if !resp.Found {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"found": true,
		"user":  resp.User,
	})
}

func listUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.UserClient.ListUsers(ctx, &proto.ListUsersRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"users": resp.Users,
		"total": resp.Total,
		"page":  resp.Page,
		"limit": resp.Limit,
	})
}

func updateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int32  `json:"age"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.UserClient.UpdateUser(ctx, &proto.UpdateUserRequest{
		UserId: int32(id),
		Name:   req.Name,
		Email:  req.Email,
		Age:    req.Age,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": resp.Success,
		"message": resp.Message,
		"user":    resp.User,
	})
}

func deleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.UserClient.DeleteUser(ctx, &proto.DeleteUserRequest{
		UserId: int32(id),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": resp.Success,
		"message": resp.Message,
	})
}

// Product endpoint handlers
func createProduct(c *fiber.Ctx) error {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		UserId      int32   `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.ProductClient.CreateProduct(ctx, &proto.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		UserId:      req.UserId,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": resp.Success,
		"message": resp.Message,
		"product": resp.Product,
	})
}

func getProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.ProductClient.GetProduct(ctx, &proto.GetProductRequest{
		ProductId: int32(id),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if !resp.Found {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.JSON(fiber.Map{
		"found":   true,
		"product": resp.Product,
	})
}

func listProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.ProductClient.ListProducts(ctx, &proto.ListProductsRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"products": resp.Products,
		"total":    resp.Total,
		"page":     resp.Page,
		"limit":    resp.Limit,
	})
}

func getUserProducts(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.ProductClient.GetProductsByUser(ctx, &proto.GetProductsByUserRequest{
		UserId: int32(id),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"products": resp.Products,
		"total":    resp.Total,
	})
}