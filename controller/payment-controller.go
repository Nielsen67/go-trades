package controller

import (
	"go-trades/entity"
	"go-trades/middleware"
	"go-trades/service"
	"go-trades/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	Service service.PaymentService
}

func NewPaymentController(s service.PaymentService) *PaymentController {
	return &PaymentController{
		Service: s,
	}
}

func (c *PaymentController) GetAllPayments(ctx *gin.Context) {
	var resp *utils.Response
	var totalSize, totalPage int64
	var err error

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

	isAdmin, err := middleware.IsAdmin(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if isAdmin {
		resp, totalSize, totalPage, err = c.Service.GetAllPayments(ctx, page, size)
	} else {
		resp, totalSize, totalPage, err = c.Service.GetUserPayments(ctx, userId.(uint), page, size)
	}

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("x-total-count", strconv.FormatInt(totalSize, 10))
	ctx.Header("x-total-page", strconv.FormatInt(totalPage, 10))

	ctx.JSON(200, resp)
}

func (c *PaymentController) CreatePayment(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(400, gin.H{"error": "user ID not found in context"})
		return
	}

	var req entity.PaymentRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}
	resp, err := c.Service.CreatePayment(ctx, userId.(uint), &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, resp)
}
