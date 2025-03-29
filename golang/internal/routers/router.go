package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntthanh2603/10day-golang.git/golang/internal/controller"
	middlewares "github.com/ntthanh2603/10day-golang.git/golang/internal/midderwares"
)

// func Middleware_A() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("Before --> Middleware A")
// 		c.Next()
// 		fmt.Println("After --> Middleware A")
// 	}
// }
// func Middleware_B(c *gin.Context) {
// 	fmt.Println("Before --> Middleware B")
// 	c.Next()
// 	fmt.Println("After --> Middleware B")
// }

func NewRouter() *gin.Engine {
	r := gin.Default()

	// use middleware
	r.Use(middlewares.AuthMiddleware())

	// Version 1
	v1 := r.Group("/v1")
	{
		v1.GET("/user", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "User page",
			})
		})
		v1.GET("/product", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Product page",
			})
		})
	}

	// Home page not version
	r.GET("/home", controller.NewUserController().GetUserByID)

	return r
}
