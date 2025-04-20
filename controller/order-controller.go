package controller

import (
	"go-trades/entity"
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
	statusStr := ctx.Query("status")
	var status uint

	if statusStr != "" {
		parsedId, err := strconv.ParseUint(statusStr, 10, 32)
		if err != nil {
			ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidOrderStatus})
			return
		}
		status = uint(parsedId)
	}

	resp, err := c.Service.GetAllOrders(ctx, status)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *OrderController) GetOrderById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
		return
	}

	resp, err := c.Service.GetOrderById(ctx, uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req entity.CreateOrderRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}
	resp, err := c.Service.CreateOrder(ctx, &req)
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
		return
	}
	resp, err := c.Service.ConfirmOrder(ctx, uint(id))
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *OrderController) CancelOrder(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidProductId})
		return
	}

	if err := c.Service.CancelOrder(ctx, uint(id)); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(204, nil)
}
