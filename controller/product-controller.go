package controller

import (
	"go-trades/entity"
	"go-trades/service"
	"go-trades/utils"
	"strconv"

	errorMessages "go-trades/utils/error-messages"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	Service service.ProductService
}

func NewProductController(s service.ProductService) *ProductController {
	return &ProductController{
		Service: s,
	}
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {

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

	categoryIdStr := ctx.Query("categoryId")
	var categoryId uint

	if categoryIdStr != "" {
		parsedId, err := strconv.ParseUint(categoryIdStr, 10, 32)
		if err != nil {
			ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidCategoryId})
			return
		}
		categoryId = uint(parsedId)
	}

	resp, totalSize, totalPage, err := c.Service.GetAllProducts(ctx, page, size, categoryId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("x-total-count", strconv.FormatInt(totalSize, 10))
	ctx.Header("x-total-page", strconv.FormatInt(totalPage, 10))

	ctx.JSON(200, resp)
}

func (c *ProductController) GetProductById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidProductId})
		return
	}

	resp, err := c.Service.GetProductById(ctx, uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req entity.CreateProductRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}
	resp, err := c.Service.CreateProduct(ctx, &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, resp)
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	var req entity.UpdateProductRequest
	if err := utils.ValidateJson(ctx, &req); err != nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidProductId})
		return
	}

	resp, err := c.Service.UpdateProduct(ctx, uint(id), &req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidProductId})
		return
	}

	if err := c.Service.DeleteProduct(ctx, uint(id)); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(204, nil)
}
