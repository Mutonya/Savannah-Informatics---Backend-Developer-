package main

import (
	"context"
	"fmt"
	"github.com/Mutonya/Savanah/internal/domain/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Mutonya/Savanah/internal/config"
	"github.com/Mutonya/Savanah/internal/controllers"
	"github.com/Mutonya/Savanah/internal/domain/repositories"
	"github.com/Mutonya/Savanah/internal/domain/services"
	"github.com/Mutonya/Savanah/internal/middleware"
	"github.com/Mutonya/Savanah/internal/routes"
	"github.com/Mutonya/Savanah/internal/utils/logging"
	"github.com/Mutonya/Savanah/pkg/database"
	"github.com/Mutonya/Savanah/pkg/oauth2"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger := logging.NewLogger(os.Stdout, cfg.Environment)

	// Initialize database
	db, err := database.NewPostgresDB(&database.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.SSLMode,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Initialize OAuth provider
	fmt.Println("OIDC Provider URL:", cfg.OAuthProviderURL)
	oauthProvider, err := oauth2.NewOIDCProvider(
		context.Background(),
		cfg.OAuthClientID,
		cfg.OAuthClientSecret,
		cfg.OAuthRedirectURL,
		cfg.OAuthProviderURL,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize OAuth provider")
	}

	// Initialize repositories
	customerRepo := repositories.NewCustomerRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Initialize services
	authService := services.NewAuthService(oauthProvider, customerRepo, cfg)
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	notificationService := services.NewNotificationService(cfg)
	orderService := services.NewOrderService(orderRepo, productRepo, customerRepo, notificationService)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	productController := controllers.NewProductController(productService)
	categoryController := controllers.NewCategoryController(categoryService)
	orderController := controllers.NewOrderController(orderService, notificationService)

	// Create Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logging.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())

	// Setup routes
	routes.SetupAuthRoutes(router, authController)
	routes.SetupAPIRoutes(router, authService, productController, categoryController, orderController, authController)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	logger.Info().Msgf("Server started on port %s", cfg.ServerPort)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited properly")
}

func runMigrations(db *gorm.DB) error {
	// Automatic migrations for simple cases
	err := db.AutoMigrate(
		&models.Customer{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return err
	}

	// Manual SQL migrations for complex changes
	// (Would use a migration tool like golang-migrate in production)
	return nil
}
