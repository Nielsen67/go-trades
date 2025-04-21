package main

import (
	"go-trades/controller"
	"go-trades/entity"
	"go-trades/middleware"
	"go-trades/repository"
	"go-trades/service"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func routeInit(conn *gorm.DB) *gin.Engine {
	// ============== Dependency Injection ============

	userRepository := repository.NewUserRepository(conn)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	categoryRepository := repository.NewCategoryRepository(conn)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)

	productRepository := repository.NewProductRepository(conn)
	productService := service.NewProductService(productRepository, categoryRepository)
	productController := controller.NewProductController(productService)

	inventoryRepository := repository.NewInventoryRepository(conn)
	inventoryService := service.NewInventoryService(inventoryRepository, productRepository)
	inventoryController := controller.NewInventoryController(inventoryService)

	orderRepository := repository.NewOrderRepository(conn)
	orderService := service.NewOrderService(conn, orderRepository, productRepository, inventoryRepository)
	orderController := controller.NewOrderController(orderService)

	paymentRepository := repository.NewPaymentRepository(conn)
	paymentService := service.NewPaymentService(conn, paymentRepository, orderRepository)
	paymentController := controller.NewPaymentController(paymentService)

	reportRepository := repository.NewReportRepository(conn)
	reportService := service.NewReportService(reportRepository)
	reportController := controller.NewReportController(reportService)

	r := gin.Default()
	api := r.Group("/api/v1")

	api.POST("/register", userController.Register)
	api.POST("/login", userController.Login)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes (Admin & Customer)
		protected.PUT("/user/password", userController.ChangePassword)
		protected.GET("/user/me", userController.GetUser)

		// Routes accessible to both Admin & Customer
		bothRoles := protected.Group("")
		bothRoles.Use(middleware.RBACMiddleware(entity.Admin, entity.Customer))
		{
			bothRoles.GET("/products", productController.GetAllProducts)
			bothRoles.GET("/products/:id", productController.GetProductById)
			bothRoles.GET("/orders", orderController.GetAllOrders)
			bothRoles.GET("/orders/:id", orderController.GetOrderById)
			bothRoles.GET("/payments", paymentController.GetAllPayments)
		}

		// Admin-only routes
		admin := protected.Group("")
		admin.Use(middleware.RBACMiddleware(entity.Admin))
		{
			// Category routes
			admin.GET("/categories", categoryController.GetAllCategories)
			admin.GET("/categories/:id", categoryController.GetCategoryById)
			admin.POST("/categories", categoryController.CreateCategory)
			admin.PUT("/categories/:id", categoryController.UpdateCategory)
			admin.DELETE("/categories/:id", categoryController.DeleteCategory)

			// Product routes (write operations)
			admin.POST("/products", productController.CreateProduct)
			admin.PUT("/products/:id", productController.UpdateProduct)
			admin.DELETE("/products/:id", productController.DeleteProduct)

			// Inventory routes
			admin.GET("/inventories", inventoryController.GetAllInventories)
			admin.GET("/inventories/:id", inventoryController.GetInventoryById)
			admin.POST("/inventories", inventoryController.CreateInventory)
			admin.PUT("/inventories/:id", inventoryController.UpdateInventory)
			admin.DELETE("/inventories/:id", inventoryController.DeleteInventory)

			// Order routes
			admin.POST("/orders/:id/process", orderController.ProcessOrder)

			// Report routes
			admin.GET("/reports", reportController.GetReport)
		}

		// Customer-only routes
		customer := protected.Group("")
		customer.Use(middleware.RBACMiddleware(entity.Customer))
		{
			// Order routes
			customer.POST("/orders", orderController.CreateOrder)
			customer.POST("/orders/:id/confirm", orderController.ConfirmOrder)
			customer.DELETE("/orders/:id", orderController.CancelOrder)

			// Payment routes
			customer.POST("/payments", paymentController.CreatePayment)
		}
	}

	return r
}
