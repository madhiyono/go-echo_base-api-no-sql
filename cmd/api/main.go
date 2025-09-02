package main

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/config"
	"github.com/madhiyono/base-api-nosql/internal/auth"
	"github.com/madhiyono/base-api-nosql/internal/handlers"
	"github.com/madhiyono/base-api-nosql/internal/middleware"
	mongorepo "github.com/madhiyono/base-api-nosql/internal/repository/mongo"
	"github.com/madhiyono/base-api-nosql/internal/routes"
	"github.com/madhiyono/base-api-nosql/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Initialize Configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error Loading Config: %v", err)
	}

	// Initialize Logger
	logger := logger.New(cfg.LogLevel)

	// Initialize MongoDB Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURL))
	if err != nil {
		logger.Fatal("Failed to Connect To MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logger.Fatal("Fatal to Disconnect From MongoDB: %v", err)
		}
	}()

	// Check Connection
	if err := client.Ping(ctx, nil); err != nil {
		logger.Fatal("Failed to Ping MongoDB: %v", err)
	}

	db := client.Database(cfg.DatabaseName)

	// Initialize Repositories
	userRepo := mongorepo.NewUserRepository(db)
	authRepo := mongorepo.NewAuthRepository(db)

	// Initialize Auth Service & Middleware
	authService := auth.NewAuthService(authRepo, userRepo, cfg.JWTSecret)
	authMiddleware := auth.NewMiddleware(authService)

	// Initialize Handlers
	userHandler := handlers.NewUserHandler(userRepo, logger)
	authHandler := handlers.NewAuthHandler(authService, logger)

	// Initialize Echo Instance
	e := echo.New()

	// Initialize Middleware
	middleware.Init(e, logger)

	// Setup Routes
	routes.Setup(e, userHandler, authHandler, authMiddleware)

	// Start Server
	logger.Info("Starting Server on Port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to Start Server: %v", err)
	}
}
