package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/auth"
	"github.com/madhiyono/base-api-nosql/internal/handlers"
	"github.com/madhiyono/base-api-nosql/internal/models"
)

func Setup(
	e *echo.Echo,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	roleHandler *handlers.RoleHandler,
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

	// Role Management Routes (Admin Only)
	roleRoutes := protected.Group("/roles")
	roleRoutes.Use(authMiddleware.RequireAdmin())
	{
		roleRoutes.POST("", roleHandler.CreateRole)
		roleRoutes.GET("/:id", roleHandler.GetRole)
		roleRoutes.PUT("/:id", roleHandler.UpdateRole)
		roleRoutes.DELETE("/:id", roleHandler.DeleteRole)
		roleRoutes.GET("", roleHandler.ListRoles)
	}

	// User Routes (Authenticated Users)
	userRoutes := protected.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser, authMiddleware.RequirePermission(models.ResourceUsers, models.ActionCreate))
		userRoutes.GET("/:id", userHandler.GetUser, authMiddleware.RequirePermission(models.ResourceUsers, models.ActionRead))
		userRoutes.PUT("/:id", userHandler.UpdateUser, authMiddleware.RequirePermission(models.ResourceUsers, models.ActionUpdate))
		userRoutes.DELETE("/:id", userHandler.DeleteUser, authMiddleware.RequirePermission(models.ResourceUsers, models.ActionDelete))
		userRoutes.GET("", userHandler.ListUsers, authMiddleware.RequirePermission(models.ResourceUsers, models.ActionRead))
	}

	// Admin Only Example
	adminRoutes := protected.Group("/admin")
	adminRoutes.Use(authMiddleware.RequireAdmin())
	{
		adminRoutes.GET("", func(c echo.Context) error {
			return c.JSON(200, "Admin Access Only!")
		})
	}
}
