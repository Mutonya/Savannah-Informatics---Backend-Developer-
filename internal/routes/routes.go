package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Mutonya/Savanah/internal/controllers"
	"github.com/Mutonya/Savanah/internal/domain/services"
	"github.com/Mutonya/Savanah/internal/middleware"
)

func SetupHealthRoute(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

func SetupAuthRoutes(router *gin.Engine, authController *controllers.AuthController) {
	auth := router.Group("/auth")
	{
		auth.GET("/login", authController.Login)
		auth.GET("/callback", authController.Callback)
	}
}

func SetupAPIRoutes(
	router *gin.Engine,
	authService services.AuthService,
	productController *controllers.ProductController,
	categoryController *controllers.CategoryController,
	orderController *controllers.OrderController,
	authController *controllers.AuthController,
) {
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(authService))
	{
		// Customer routes
		api.GET("/profile", authController.Profile)

		// Product routes
		api.POST("/products", productController.CreateProduct)
		api.GET("/products", productController.GetProducts)
		api.GET("/products/:id", productController.GetProduct)
		api.PUT("/products/:id", productController.UpdateProduct)
		api.DELETE("/products/:id", productController.DeleteProduct)

		// Category routes
		api.POST("/categories", categoryController.CreateCategory)
		api.GET("/categories", categoryController.GetCategories)
		api.GET("/categories/:id", categoryController.GetCategory)
		api.GET("/categories/:id/products", categoryController.GetCategoryProducts)
		api.GET("/categories/:id/average-price", categoryController.GetAveragePrice)
		api.PUT("/categories/:id", categoryController.UpdateCategory)
		api.DELETE("/categories/:id", categoryController.DeleteCategory)

		// Order routes
		api.POST("/orders", orderController.CreateOrder)
		api.GET("/orders", orderController.GetOrders)
		api.GET("/orders/:id", orderController.GetOrder)
		api.PUT("/orders/:id/status", orderController.UpdateOrderStatus)
	}
}
