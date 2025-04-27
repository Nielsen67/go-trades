package controller

import (
	"go-trades/entity"
	"go-trades/service"
	"go-trades/utils"
	errorMessages "go-trades/utils/error-messages"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service service.UserService
}

func NewUserController(s service.UserService) *UserController {
	return &UserController{
		Service: s,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req entity.UserRegisterRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}

	resp, err := c.Service.Register(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, resp)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req entity.UserLoginRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}

	resp, err := c.Service.Login(ctx, &req)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req entity.UserChangePassword
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}

	userIdAny, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userId, ok := userIdAny.(uint)
	if !ok {
		ctx.JSON(401, gin.H{"error": "invalid user ID type"})
		return
	}

	resp, err := c.Service.ChangePassword(ctx, userId, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	userIdAny, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userId, ok := userIdAny.(uint)
	if !ok {
		ctx.JSON(401, gin.H{"error": "invalid user ID type"})
		return
	}

	resp, err := c.Service.GetUserById(ctx, userId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *UserController) AssignAsAdmin(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrUserNotExists})
		return
	}

	err = c.Service.AssignAsAdmin(ctx, uint(id))
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User assigned as admin successfully"})
}
