package controller

import (
	"go-trades/entity"
	"go-trades/service"
	"go-trades/utils"
	"strconv"

	errorMessages "go-trades/utils/error-messages"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	Service service.CategoryService
}

func NewCategoryController(s service.CategoryService) *CategoryController {
	return &CategoryController{
		Service: s,
	}
}

func (c *CategoryController) GetAllCategories(ctx *gin.Context) {

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

	resp, totalSize, totalPage, err := c.Service.GetAllCategories(ctx, page, size)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("x-total-count", strconv.FormatInt(totalSize, 10))
	ctx.Header("x-total-page", strconv.FormatInt(totalPage, 10))

	ctx.JSON(200, resp)
}

func (c *CategoryController) GetCategoryById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
		return
	}

	resp, err := c.Service.GetCategoryById(ctx, uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req entity.CategoryRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}
	resp, err := c.Service.CreateCategory(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, resp)
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	var req entity.CategoryRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
		return
	}

	resp, err := c.Service.UpdateCategory(ctx, uint(id), &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
		return
	}

	if err := c.Service.DeleteCategory(ctx, uint(id)); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(204, nil)
}
