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