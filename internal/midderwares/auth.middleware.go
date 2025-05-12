package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token := c.GetHeader("authentication")
		// if token != "123456" {
		// 	response.ErrorResponse(c, response.ErrInvalidToken, "")
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}
