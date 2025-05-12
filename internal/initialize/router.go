package initialize

import (
	"net/http"

	"myproject/internal/controller"
	middlewares "myproject/internal/midderwares"

	"github.com/gin-gonic/gin"
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

func InitRouter() *gin.Engine {
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
