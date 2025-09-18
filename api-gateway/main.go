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
	UserClient      proto.UserServiceClient
	ProductClient   proto.ProductServiceClient
	InventoryClient proto.InventoryServiceClient
	OrderClient     proto.OrderServiceClient
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
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400,
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

	// Inventory routes
	inventoryRoutes := api.Group("/inventory")
	inventoryRoutes.Post("/", createInventoryItem)
	inventoryRoutes.Get("/:id", getInventoryItem)
	inventoryRoutes.Put("/:id", updateInventoryItem)
	inventoryRoutes.Get("/", listInventoryItems)
	inventoryRoutes.Post("/check-stock", checkStock)
	inventoryRoutes.Post("/reserve-stock", reserveStock)
	inventoryRoutes.Post("/release-stock", releaseStock)

	// Order routes
	orderRoutes := api.Group("/orders")
	orderRoutes.Post("/", createOrder)
	orderRoutes.Get("/:id", getOrder)
	orderRoutes.Get("/", listOrders)
	orderRoutes.Put("/:id/status", updateOrderStatus)

	log.Println("🚀 API Gateway starting on port 8000")
	log.Println("📍 User endpoints: /api/users")
	log.Println("📍 Product endpoints: /api/products")
	log.Println("📍 Inventory endpoints: /api/inventory")
	log.Println("📍 Order endpoints: /api/orders")
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

	// Connect to Inventory Service
	inventoryConn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcClients{
		UserClient:      proto.NewUserServiceClient(userConn),
		ProductClient:   proto.NewProductServiceClient(productConn),
		InventoryClient: proto.NewInventoryServiceClient(inventoryConn),
		OrderClient:     proto.NewOrderServiceClient(inventoryConn),
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

// Inventory handlers

// createInventoryItem Create Inventory Item
// @Summary      Create a new inventory item
// @Description  Create a new inventory item for a product
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        item  body      models.CreateInventoryItemRequest  true  "Inventory item data"
// @Success      201   {object}  models.InventoryItemResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /inventory [post]
func createInventoryItem(c *fiber.Ctx) error {
	var req struct {
		ProductID int32  `json:"product_id"`
		Quantity  int32  `json:"quantity"`
		Location  string `json:"location"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.InventoryClient.CreateInventoryItem(ctx, &proto.CreateInventoryItemRequest{
		ProductId: req.ProductID,
		Quantity:  req.Quantity,
		Location:  req.Location,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": resp.Message,
		"item":    resp.Item,
	})
}

// getInventoryItem Get Inventory Item
// @Summary      Get inventory item by ID
// @Description  Get an inventory item by its ID
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Inventory Item ID"
// @Success      200  {object}  models.InventoryItemResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /inventory/{id} [get]
func getInventoryItem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid inventory item ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.InventoryClient.GetInventoryItem(ctx, &proto.GetInventoryItemRequest{
		Id: int32(id),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": resp.Message,
		"item":    resp.Item,
	})
}

// checkStock Check Stock
// @Summary      Check stock availability
// @Description  Check if sufficient stock is available for a product
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        request  body      models.CheckStockRequest  true  "Stock check data"
// @Success      200      {object}  models.CheckStockResponse
// @Failure      400      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /inventory/check-stock [post]
func checkStock(c *fiber.Ctx) error {
	var req struct {
		ProductID        int32 `json:"product_id"`
		RequiredQuantity int32 `json:"required_quantity"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.InventoryClient.CheckStock(ctx, &proto.CheckStockRequest{
		ProductId:        req.ProductID,
		RequiredQuantity: req.RequiredQuantity,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"available":           resp.Available,
		"available_quantity": resp.AvailableQuantity,
		"message":            resp.Message,
	})
}

// reserveStock Reserve Stock
// @Summary      Reserve stock for an order
// @Description  Reserve stock for a specific order
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        request  body      models.ReserveStockRequest  true  "Stock reservation data"
// @Success      200      {object}  models.ReserveStockResponse
// @Failure      400      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /inventory/reserve-stock [post]
func reserveStock(c *fiber.Ctx) error {
	var req struct {
		ProductID int32  `json:"product_id"`
		Quantity  int32  `json:"quantity"`
		OrderID   string `json:"order_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.InventoryClient.ReserveStock(ctx, &proto.ReserveStockRequest{
		ProductId: req.ProductID,
		Quantity:  req.Quantity,
		OrderId:   req.OrderID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":        resp.Success,
		"message":        resp.Message,
		"reservation_id": resp.ReservationId,
	})
}

// releaseStock Release Stock
// @Summary      Release reserved stock
// @Description  Release stock that was previously reserved
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        request  body      models.ReleaseStockRequest  true  "Stock release data"
// @Success      200      {object}  models.ReleaseStockResponse
// @Failure      400      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /inventory/release-stock [post]
func releaseStock(c *fiber.Ctx) error {
	var req struct {
		ReservationID string `json:"reservation_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.InventoryClient.ReleaseStock(ctx, &proto.ReleaseStockRequest{
		ReservationId: req.ReservationID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": resp.Success,
		"message": resp.Message,
	})
}

// Order handlers

// createOrder Create Order
// @Summary      Create a new order
// @Description  Create a new order with items
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order  body      models.CreateOrderRequest  true  "Order data"
// @Success      201    {object}  models.OrderResponse
// @Failure      400    {object}  models.ErrorResponse
// @Failure      500    {object}  models.ErrorResponse
// @Router       /orders [post]
func createOrder(c *fiber.Ctx) error {
	var req struct {
		UserID int32 `json:"user_id"`
		Items  []struct {
			ProductID int32   `json:"product_id"`
			Quantity  int32   `json:"quantity"`
			Price     float64 `json:"price"`
		} `json:"items"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert items
	var orderItems []*proto.OrderItem
	for _, item := range req.Items {
		orderItems = append(orderItems, &proto.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	resp, err := clients.OrderClient.CreateOrder(ctx, &proto.CreateOrderRequest{
		UserId: req.UserID,
		Items:  orderItems,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": resp.Message,
		"order":   resp.Order,
	})
}

// getOrder Get Order
// @Summary      Get order by ID
// @Description  Get an order by its ID
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  models.OrderResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /orders/{id} [get]
func getOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.OrderClient.GetOrder(ctx, &proto.GetOrderRequest{
		Id: id,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": resp.Message,
		"order":   resp.Order,
	})
}

// listOrders List Orders
// @Summary      List orders with pagination
// @Description  Get a paginated list of orders, optionally filtered by user
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        user_id  query     int  false  "User ID filter"
// @Param        page     query     int  false  "Page number"  default(1)
// @Param        limit    query     int  false  "Items per page"  default(10)
// @Success      200      {object}  models.OrdersListResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /orders [get]
func listOrders(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Query("user_id", "0"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.OrderClient.ListOrders(ctx, &proto.ListOrdersRequest{
		UserId: int32(userID),
		Page:   int32(page),
		Limit:  int32(limit),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"orders": resp.Orders,
		"total":  resp.Total,
		"page":   resp.Page,
		"limit":  resp.Limit,
	})
}

// updateOrderStatus Update Order Status
// @Summary      Update order status
// @Description  Update the status of an order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id      path      string                      true  "Order ID"
// @Param        status  body      models.UpdateOrderStatusRequest  true  "Status update data"
// @Success      200     {object}  models.OrderResponse
// @Failure      400     {object}  models.ErrorResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /orders/{id}/status [put]
func updateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Convert string status to enum
	var status proto.OrderStatus
	switch req.Status {
	case "PENDING":
		status = proto.OrderStatus_PENDING
	case "CONFIRMED":
		status = proto.OrderStatus_CONFIRMED
	case "PROCESSING":
		status = proto.OrderStatus_PROCESSING
	case "SHIPPED":
		status = proto.OrderStatus_SHIPPED
	case "DELIVERED":
		status = proto.OrderStatus_DELIVERED
	case "CANCELLED":
		status = proto.OrderStatus_CANCELLED
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.OrderClient.UpdateOrderStatus(ctx, &proto.UpdateOrderStatusRequest{
		Id:     id,
		Status: status,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": resp.Message,
		"order":   resp.Order,
	})
}

// listInventoryItems List Inventory Items
// @Summary      List inventory items with pagination
// @Description  Get a paginated list of all inventory items
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"  default(1)
// @Param        limit  query     int  false  "Items per page"  default(10)
// @Success      200    {object}  models.InventoryItemsListResponse
// @Failure      500    {object}  models.ErrorResponse
// @Router       /inventory [get]
func listInventoryItems(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.InventoryClient.ListInventoryItems(ctx, &proto.ListInventoryItemsRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"items": resp.Items,
		"total": resp.Total,
		"page":  resp.Page,
		"limit": resp.Limit,
	})
}

// updateInventoryItem Update Inventory Item
// @Summary      Update inventory item
// @Description  Update an inventory item's quantity and location
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        id    path      int                           true  "Inventory Item ID"
// @Param        item  body      models.UpdateInventoryItemRequest  true  "Inventory item update data"
// @Success      200   {object}  models.InventoryItemResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /inventory/{id} [put]
func updateInventoryItem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid inventory item ID"})
	}

	var req struct {
		Quantity int32  `json:"quantity"`
		Location string `json:"location"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := clients.InventoryClient.UpdateInventoryItem(ctx, &proto.UpdateInventoryItemRequest{
		Id:       int32(id),
		Quantity: req.Quantity,
		Location: req.Location,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": resp.Message,
		"item":    resp.Item,
	})
}