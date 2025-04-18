package main

import (
	"go-trades/controller"
	"go-trades/repository"
	"go-trades/service"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func routeInit(conn *gorm.DB) *gin.Engine {
	// ============== Dependency Injection ============
	categoryRepository := repository.NewCategoryRepository(conn)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)

	productRepository := repository.NewProductRepository(conn)
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	inventoryRepository := repository.NewInventoryRepository(conn)
	inventoryService := service.NewInventoryService(inventoryRepository)
	inventoryController := controller.NewInventoryController(inventoryService)

	r := gin.Default()
	api := r.Group("/api/v1")
	{
		//CATEGORY ROUTES
		api.GET("/categories", categoryController.GetAllCategories)
		api.GET("/categories/:id", categoryController.GetCategoryById)
		api.POST("/categories", categoryController.CreateCategory)
		api.PUT("/categories/:id", categoryController.UpdateCategory)
		api.DELETE("/categories/:id", categoryController.DeleteCategory)

		//PRODUCT ROUTES
		api.GET("/products", productController.GetAllProducts)
		api.GET("/products/:id", productController.GetProductById)
		api.POST("/products", productController.CreateProduct)
		api.PUT("/products/:id", productController.UpdateProduct)
		api.DELETE("/products/:id", productController.DeleteProduct)

		//INVENTORY ROUTES
		api.GET("/inventories", inventoryController.GetAllInventories)
		api.GET("/inventories/:id", inventoryController.GetInventoryById)
		api.POST("/inventories", inventoryController.CreateInventory)
		api.PUT("/inventories/:id", inventoryController.UpdateInventory)
		api.DELETE("/inventories/:id", inventoryController.DeleteInventory)
	}
	return r
}
