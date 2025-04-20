package controller

import (
	"go-trades/entity"
	"go-trades/service"
	"go-trades/utils"

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

	resp, err := c.Service.GetAllPayments(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *PaymentController) CreatePayment(ctx *gin.Context) {
	var req entity.PaymentRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}
	resp, err := c.Service.CreatePayment(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, resp)
}
