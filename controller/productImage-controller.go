package controller

import (
	"go-trades/service"
	errorMessages "go-trades/utils/error-messages"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductImageController struct {
	Service service.ProductImageService
}

func NewProductImageController(s service.ProductImageService) *ProductImageController {
	return &ProductImageController{
		Service: s,
	}
}

func (c *ProductImageController) UploadProductImage(ctx *gin.Context) {
	productId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": errorMessages.ErrInvalidProductId})
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Failed to upload image"})
		return
	}

	err = c.Service.UploadProductImage(ctx, uint(productId), image)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Image uploaded successfully",
	})
}

func (c *ProductImageController) DownloadProductImages(ctx *gin.Context) {
	rawProductId := ctx.Param("id")
	productId, err := strconv.ParseUint(rawProductId, 10, 0)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	zipFileName, err := c.Service.DownloadProductImages(ctx, uint(productId))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", "attachment; filename="+filepath.Base(zipFileName))
	ctx.File(zipFileName)
}
