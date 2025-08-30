package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/handlers"
)

func Setup(e *echo.Echo, userHandler *handlers.UserHandler) {
	// Root Endpoint
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World!")
	})

	// Health Check Endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

	// Register Routes
	userHandler.RegisterRoutes(e)
}