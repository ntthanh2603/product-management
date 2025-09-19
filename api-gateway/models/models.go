package models

import "time"

// User represents a user in the system
// @Description User information
type User struct {
	ID        int32     `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	Age       int32     `json:"age" example:"30"`
	CreatedAt string    `json:"created_at" example:"2023-01-01T12:00:00Z"`
} //@name User

// CreateUserRequest request to create a new user
// @Description Request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required" example:"john@example.com"`
	Age   int32  `json:"age" binding:"required" example:"30"`
} //@name CreateUserRequest

// UpdateUserRequest request to update an existing user
// @Description Request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required" example:"john@example.com"`
	Age   int32  `json:"age" binding:"required" example:"30"`
} //@name UpdateUserRequest

// Product represents a product in the system
// @Description Product information
type Product struct {
	ID          int32     `json:"id" example:"1"`
	Name        string    `json:"name" example:"iPhone 15"`
	Description string    `json:"description" example:"Latest iPhone model"`
	Price       float64   `json:"price" example:"999.99"`
	UserID      int32     `json:"user_id" example:"1"`
	CreatedAt   string    `json:"created_at" example:"2023-01-01T12:00:00Z"`
} //@name Product

// CreateProductRequest request to create a new product
// @Description Request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required" example:"iPhone 15"`
	Description string  `json:"description" binding:"required" example:"Latest iPhone model"`
	Price       float64 `json:"price" binding:"required" example:"999.99"`
	UserID      int32   `json:"user_id" binding:"required" example:"1"`
} //@name CreateProductRequest

// SuccessResponse represents a successful operation response
// @Description Success response
type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
} //@name SuccessResponse

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
	Code  int    `json:"code" example:"400"`
} //@name ErrorResponse

// UserResponse represents a user response
// @Description User response
type UserResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"User created successfully"`
	User    User   `json:"user"`
} //@name UserResponse

// UsersListResponse represents a list of users response
// @Description Users list response
type UsersListResponse struct {
	Users []User `json:"users"`
	Total int32  `json:"total" example:"100"`
	Page  int32  `json:"page" example:"1"`
	Limit int32  `json:"limit" example:"10"`
} //@name UsersListResponse

// ProductResponse represents a product response
// @Description Product response
type ProductResponse struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Product created successfully"`
	Product Product `json:"product"`
} //@name ProductResponse

// ProductsListResponse represents a list of products response
// @Description Products list response
type ProductsListResponse struct {
	Products []Product `json:"products"`
	Total    int32     `json:"total" example:"50"`
	Page     int32     `json:"page" example:"1"`
	Limit    int32     `json:"limit" example:"10"`
} //@name ProductsListResponse

// UserProductsResponse represents products of a specific user
// @Description User products response
type UserProductsResponse struct {
	Products []Product `json:"products"`
	Total    int32     `json:"total" example:"5"`
} //@name UserProductsResponse

// HealthResponse represents health check response
// @Description Health check response
type HealthResponse struct {
	Status   string                 `json:"status" example:"healthy"`
	Time     time.Time              `json:"time" example:"2023-01-01T12:00:00Z"`
	Services map[string]string      `json:"services"`
} //@name HealthResponse

// InventoryItem represents an inventory item
// @Description Inventory item information
type InventoryItem struct {
	ID        int32  `json:"id" example:"1"`
	ProductID int32  `json:"product_id" example:"1"`
	Quantity  int32  `json:"quantity" example:"100"`
	Location  string `json:"location" example:"Warehouse A"`
	CreatedAt string `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01T12:00:00Z"`
} //@name InventoryItem

// CreateInventoryItemRequest request to create a new inventory item
// @Description Request body for creating an inventory item
type CreateInventoryItemRequest struct {
	ProductID int32  `json:"product_id" binding:"required" example:"1"`
	Quantity  int32  `json:"quantity" binding:"required" example:"100"`
	Location  string `json:"location" binding:"required" example:"Warehouse A"`
} //@name CreateInventoryItemRequest

// UpdateInventoryItemRequest request to update an inventory item
// @Description Request body for updating an inventory item
type UpdateInventoryItemRequest struct {
	Quantity int32  `json:"quantity" binding:"required" example:"100"`
	Location string `json:"location" binding:"required" example:"Warehouse A"`
} //@name UpdateInventoryItemRequest

// InventoryItemResponse represents an inventory item response
// @Description Inventory item response
type InventoryItemResponse struct {
	Message string        `json:"message" example:"Inventory item created successfully"`
	Item    InventoryItem `json:"item"`
} //@name InventoryItemResponse

// InventoryItemsListResponse represents a list of inventory items response
// @Description Inventory items list response
type InventoryItemsListResponse struct {
	Items []InventoryItem `json:"items"`
	Total int32           `json:"total" example:"50"`
	Page  int32           `json:"page" example:"1"`
	Limit int32           `json:"limit" example:"10"`
} //@name InventoryItemsListResponse

// CheckStockRequest request to check stock availability
// @Description Request body for checking stock availability
type CheckStockRequest struct {
	ProductID        int32 `json:"product_id" binding:"required" example:"1"`
	RequiredQuantity int32 `json:"required_quantity" binding:"required" example:"10"`
} //@name CheckStockRequest

// CheckStockResponse represents stock availability response
// @Description Stock availability response
type CheckStockResponse struct {
	Available         bool   `json:"available" example:"true"`
	AvailableQuantity int32  `json:"available_quantity" example:"100"`
	Message           string `json:"message" example:"Stock is available"`
} //@name CheckStockResponse

// ReserveStockRequest request to reserve stock
// @Description Request body for reserving stock
type ReserveStockRequest struct {
	ProductID int32  `json:"product_id" binding:"required" example:"1"`
	Quantity  int32  `json:"quantity" binding:"required" example:"10"`
	OrderID   string `json:"order_id" binding:"required" example:"ord_123456"`
} //@name ReserveStockRequest

// ReserveStockResponse represents stock reservation response
// @Description Stock reservation response
type ReserveStockResponse struct {
	Success       bool   `json:"success" example:"true"`
	Message       string `json:"message" example:"Stock reserved successfully"`
	ReservationID string `json:"reservation_id" example:"res_123456"`
} //@name ReserveStockResponse

// ReleaseStockRequest request to release stock
// @Description Request body for releasing stock
type ReleaseStockRequest struct {
	ReservationID string `json:"reservation_id" binding:"required" example:"res_123456"`
} //@name ReleaseStockRequest

// ReleaseStockResponse represents stock release response
// @Description Stock release response
type ReleaseStockResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Stock released successfully"`
} //@name ReleaseStockResponse

// OrderItem represents an item in an order
// @Description Order item information
type OrderItem struct {
	ProductID int32   `json:"product_id" example:"1"`
	Quantity  int32   `json:"quantity" example:"2"`
	Price     float64 `json:"price" example:"999.99"`
} //@name OrderItem

// Order represents an order in the system
// @Description Order information
type Order struct {
	ID         string      `json:"id" example:"ord_123456"`
	UserID     int32       `json:"user_id" example:"1"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price" example:"1999.98"`
	Status     string      `json:"status" example:"PENDING"`
	CreatedAt  string      `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt  string      `json:"updated_at" example:"2023-01-01T12:00:00Z"`
} //@name Order

// CreateOrderRequest request to create a new order
// @Description Request body for creating an order
type CreateOrderRequest struct {
	UserID int32       `json:"user_id" binding:"required" example:"1"`
	Items  []OrderItem `json:"items" binding:"required"`
} //@name CreateOrderRequest

// OrderResponse represents an order response
// @Description Order response
type OrderResponse struct {
	Message string `json:"message" example:"Order created successfully"`
	Order   Order  `json:"order"`
} //@name OrderResponse

// OrdersListResponse represents a list of orders response
// @Description Orders list response
type OrdersListResponse struct {
	Orders []Order `json:"orders"`
	Total  int32   `json:"total" example:"25"`
	Page   int32   `json:"page" example:"1"`
	Limit  int32   `json:"limit" example:"10"`
} //@name OrdersListResponse

// UpdateOrderStatusRequest request to update order status
// @Description Request body for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required" example:"CONFIRMED"`
} //@name UpdateOrderStatusRequest