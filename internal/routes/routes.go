package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/auth"
	"github.com/madhiyono/base-api-nosql/internal/handlers"
)

func Setup(
	e *echo.Echo,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	authMiddleware *auth.Middleware,
) {
	// Root Endpoint
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World!")
	})

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

	// Auth Routes (No Authentication Required)
	authHandler.RegisterRoutes(e)

	// Protected Routes
	protected := e.Group("")
	protected.Use(authMiddleware.JWTAuth)

	// User Routes (Authenticated Users)
	userRoutes := protected.Group("/users")
	userRoutes.Use(authMiddleware.RequireUser)
	{
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
		userRoutes.GET("", userHandler.ListUsers)
	}

	// Admin Only Example
	adminRoutes := protected.Group("/admin")
	adminRoutes.Use(authMiddleware.RequireAdmin)
	{
		adminRoutes.GET("", func(c echo.Context) error {
			return c.JSON(200, "Admin Access Only!")
		})
	}
}
