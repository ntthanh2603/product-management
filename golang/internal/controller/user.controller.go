package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ntthanh2603/10day-golang.git/golang/internal/service"
	"github.com/ntthanh2603/10day-golang.git/golang/pkg/response"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// controler -> service -> repo -> models -> database
func (uc *UserController) GetUserByID(c *gin.Context) {

	response.SuccessResponse(c, 20003, uc.userService.GetInfoUser())
}
