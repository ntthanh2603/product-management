package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ntthanh2603/10day-golang.git/golang/pkg/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authentication")
		if token != "123456" {
			response.ErrorResponse(c, response.ErrInvalidToken, "")
			c.Abort()
			return
		}
		c.Next()
	}
}
