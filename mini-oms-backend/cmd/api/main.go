package main

import (
	"log"
	"mini-oms-backend/internal/config"
	"mini-oms-backend/internal/db"
	"mini-oms-backend/internal/middlewares"
	"mini-oms-backend/internal/modules/auth"
	"mini-oms-backend/internal/modules/order"
	"mini-oms-backend/internal/modules/payment"
	"mini-oms-backend/internal/modules/product"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()
	cfg.LogConfig()

	// Connect to database
	if err := db.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize repositories
	authRepo := auth.NewRepository(db.GetDB())
	productRepo := product.NewRepository(db.GetDB())
	orderRepo := order.NewRepository(db.GetDB())
	paymentRepo := payment.NewRepository(db.GetDB())

	// Initialize services
	authService := auth.NewService(authRepo, cfg)
	productService := product.NewService(productRepo)
	orderService := order.NewService(orderRepo, db.GetDB())
	paymentService := payment.NewService(paymentRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	productHandler := product.NewHandler(productService)
	orderHandler := order.NewHandler(orderService)
	paymentHandler := payment.NewHandler(paymentService)

	// Routes
	api := e.Group("/api")

	// Public routes
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// Public product routes (anyone can view)
	api.GET("/products", productHandler.GetAll)
	api.GET("/products/:id", productHandler.GetByID)

	// Protected routes (require JWT)
	protected := api.Group("")
	protected.Use(middlewares.JWTMiddleware(cfg))

	// Order routes (protected)
	protected.GET("/orders", orderHandler.GetAll)      // User sees own, Admin sees all
	protected.GET("/orders/:id", orderHandler.GetByID) // User sees own, Admin sees all
	protected.POST("/orders", orderHandler.Create)

	// Payment routes (protected)
	protected.POST("/payments", paymentHandler.Create)
	protected.GET("/payments/order/:orderId", paymentHandler.GetByOrderID)

	// Admin-only routes
	admin := api.Group("")
	admin.Use(middlewares.JWTMiddleware(cfg))
	admin.Use(middlewares.AdminOnlyMiddleware())

	// Product management (admin only)
	admin.POST("/products", productHandler.Create)
	admin.PUT("/products/:id", productHandler.Update)
	admin.DELETE("/products/:id", productHandler.Delete)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
