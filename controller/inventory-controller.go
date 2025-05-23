package controller

import (
	"go-trades/entity"
	"go-trades/service"
	"go-trades/utils"
	"strconv"

	errorMessages "go-trades/utils/error-messages"

	"github.com/gin-gonic/gin"
)

type InventoryController struct {
	Service service.InventoryService
}

func NewInventoryController(s service.InventoryService) *InventoryController {
	return &InventoryController{
		Service: s,
	}
}

func (c *InventoryController) GetAllInventories(ctx *gin.Context) {

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

	resp, totalSize, totalPage, err := c.Service.GetAllInventories(ctx, page, size)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("x-total-count", strconv.FormatInt(totalSize, 10))
	ctx.Header("x-total-page", strconv.FormatInt(totalPage, 10))

	ctx.JSON(200, resp)
}

func (c *InventoryController) GetInventoryById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidInventoryId})
		return
	}

	resp, err := c.Service.GetInventoryById(ctx, uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *InventoryController) CreateInventory(ctx *gin.Context) {
	var req entity.CreateInventoryRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}
	resp, err := c.Service.CreateInventory(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, resp)
}

func (c *InventoryController) UpdateInventory(ctx *gin.Context) {
	var req entity.UpdateInventoryRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidInventoryId})
		return
	}

	resp, err := c.Service.UpdateInventory(ctx, uint(id), &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *InventoryController) DeleteInventory(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidInventoryId})
		return
	}

	if err := c.Service.DeleteInventory(ctx, uint(id)); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(204, nil)
}
