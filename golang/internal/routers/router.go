package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntthanh2603/go-ecommerce.git/internal/controller"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

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
