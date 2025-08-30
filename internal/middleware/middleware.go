package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/madhiyono/base-api-nosql/pkg/logger"
)

func Init(e *echo.Echo, logger *logger.Logger) {
	// Add Logger Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${host}${path} ${latency_human}\n",
	}))

	// Add Recover Middleware
	e.Use(middleware.Recover())

	// Add CORS Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Add Request ID Middleware
	e.Use(middleware.RequestID())
}