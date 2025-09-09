// Package main API Gateway for Microservices
// @title           Microservices API Gateway
// @version         1.0
// @description     API Gateway for User and Product microservices
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"api-gateway/models"
	"api-gateway/proto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "api-gateway/docs"
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

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

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
	log.Println("📖 Swagger documentation: /swagger/")

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

// healthCheck Health Check
// @Summary      Health check endpoint
// @Description  Check the health status of the API Gateway and connected services
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.HealthResponse
// @Router       /health [get]
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

// createUser Create User
// @Summary      Create a new user
// @Description  Create a new user with name, email, and age
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      models.CreateUserRequest  true  "User creation request"
// @Success      200   {object}  models.UserResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /users [post]
func createUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

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

// getUser Get User
// @Summary      Get user by ID
// @Description  Get a user by their ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{id} [get]
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

// listUsers List Users
// @Summary      List all users with pagination
// @Description  Get a paginated list of all users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"  default(1)
// @Param        limit  query     int  false  "Items per page"  default(10)
// @Success      200    {object}  models.UsersListResponse
// @Failure      500    {object}  models.ErrorResponse
// @Router       /users [get]
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

// updateUser Update User
// @Summary      Update an existing user
// @Description  Update a user's information by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      int                        true  "User ID"
// @Param        user  body      models.UpdateUserRequest  true  "User update request"
// @Success      200   {object}  models.UserResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /users/{id} [put]
func updateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req models.UpdateUserRequest

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

// deleteUser Delete User
// @Summary      Delete a user
// @Description  Delete a user by their ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{id} [delete]
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

// createProduct Create Product
// @Summary      Create a new product
// @Description  Create a new product with name, description, price and user_id
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      models.CreateProductRequest  true  "Product creation request"
// @Success      200      {object}  models.ProductResponse
// @Failure      400      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /products [post]
func createProduct(c *fiber.Ctx) error {
	var req models.CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.ProductClient.CreateProduct(ctx, &proto.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		UserId:      req.UserID,
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

// getProduct Get Product
// @Summary      Get product by ID
// @Description  Get a product by its ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /products/{id} [get]
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

// listProducts List Products
// @Summary      List all products with pagination
// @Description  Get a paginated list of all products
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"  default(1)
// @Param        limit  query     int  false  "Items per page"  default(10)
// @Success      200    {object}  models.ProductsListResponse
// @Failure      500    {object}  models.ErrorResponse
// @Router       /products [get]
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

// getUserProducts Get User Products
// @Summary      Get products by user ID
// @Description  Get all products belonging to a specific user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.UserProductsResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{id}/products [get]
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