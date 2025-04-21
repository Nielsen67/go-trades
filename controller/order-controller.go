package controller

import (
	"go-trades/entity"
	"go-trades/middleware"
	"go-trades/service"
	"go-trades/utils"
	errorMessages "go-trades/utils/error-messages"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	Service service.OrderService
}

func NewOrderController(s service.OrderService) *OrderController {
	return &OrderController{
		Service: s,
	}
}

func (c *OrderController) GetAllOrders(ctx *gin.Context) {
	var resp *utils.Response
	var totalSize, totalPage int64
	var err error
	var status uint

	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(400, gin.H{"error": "user ID not found in context"})
		return
	}

	page := utils.DefaultPage
	size := utils.DefaultSize

	var pagination utils.Pagination
	if err := ctx.ShouldBindQuery(&pagination); err == nil {
		if pagination.Page > 0 {
			page = pagination.Page
		}
		if pagination.Size > 0 {
			size = pagination.Size
		}
	}

	statusStr := ctx.Query("status")

	if statusStr != "" {
		parsedId, err := strconv.ParseUint(statusStr, 10, 32)
		if err != nil {
			ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidOrderStatus})
			return
		}
		status = uint(parsedId)
	}

	isAdmin, err := middleware.IsAdmin(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if isAdmin {
		resp, totalSize, totalPage, err = c.Service.GetAllOrders(ctx, page, size, status)
	} else {
		resp, totalSize, totalPage, err = c.Service.GetUserOrders(ctx, userId.(uint), page, size, status)
	}

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("x-total-count", strconv.FormatInt(totalSize, 10))
	ctx.Header("x-total-page", strconv.FormatInt(totalPage, 10))

	ctx.JSON(200, resp)
}

func (c *OrderController) GetOrderById(ctx *gin.Context) {
	var resp *utils.Response
	var err error

	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(400, gin.H{"error": "user ID not found in context"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
		return
	}

	isAdmin, err := middleware.IsAdmin(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if isAdmin {
		resp, err = c.Service.GetOrderById(ctx, uint(id))
	} else {
		resp, err = c.Service.GetUserOrderById(ctx, userId.(uint), uint(id))
	}

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req entity.CreateOrderRequest

	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(400, gin.H{"error": "user ID not found in context"})
		return
	}

	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}
	resp, err := c.Service.CreateOrder(ctx, userId.(uint), &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, resp)
}

func (c *OrderController) ProcessOrder(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
		return
	}
	resp, err := c.Service.ProcessOrder(ctx, uint(id))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *OrderController) ConfirmOrder(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(400, gin.H{"error": "user ID not found in context"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
		return
	}
	resp, err := c.Service.ConfirmOrder(ctx, userId.(uint), uint(id))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *OrderController) CancelOrder(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(400, gin.H{"error": "user ID not found in context"})
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidProductId})
		return
	}

	if err := c.Service.CancelOrder(ctx, userId.(uint), uint(id)); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(204, nil)
}
