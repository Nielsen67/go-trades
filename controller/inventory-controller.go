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
	resp, err := c.Service.GetAllInventories()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *InventoryController) GetInventoryById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidInventoryId})
		return
	}

	resp, err := c.Service.GetInventoryById(uint(id))
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
	resp, err := c.Service.CreateInventory(&req)
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

	resp, err := c.Service.UpdateInventory(uint(id), &req)
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

	if err := c.Service.DeleteInventory(uint(id)); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(204, nil)
}
