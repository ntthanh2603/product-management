package controller

import (
	"myproject/internal/service"

	"myproject/pkg/response"

	"github.com/gin-gonic/gin"
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
