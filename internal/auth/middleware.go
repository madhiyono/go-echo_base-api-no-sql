package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Middleware struct {
	authService *AuthService
}

func NewMiddleware(authService *AuthService) *Middleware {
	return &Middleware{
		authService: authService,
	}
}

// JWTAuth middleware for authentication
func (m *Middleware) JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, http.StatusUnauthorized, "Missing Authorization Token", nil)
		}

		// Extract token
		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			return response.Error(c, http.StatusUnauthorized, "Invalid Authorization Header Format", nil)
		}

		// Validate token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			return response.Error(c, http.StatusUnauthorized, "Invalid or Expired Token", nil)
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID)

		return next(c)
	}
}

// RequirePermission middleware for permission-based authorization
func (m *Middleware) RequirePermission(resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleID, ok := c.Get("role_id").(primitive.ObjectID)
			if !ok {
				return response.Error(c, http.StatusForbidden, "Invalid Role", nil)
			}

			// Check if user has required permission
			hasPermission, err := m.authService.HasPermission(roleID, resource, action)
			if err != nil {
				return response.InternalServerError(c, "Failed to Check Permissions", err)
			}

			if !hasPermission {
				return response.Error(c, http.StatusForbidden, "Insufficient Permissions", nil)
			}

			return next(c)
		}
	}
}

// RequireAdmin middleware for admin-only access
func (m *Middleware) RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleID, ok := c.Get("role_id").(primitive.ObjectID)
			if !ok {
				return response.Error(c, http.StatusForbidden, "Invalid Role", nil)
			}

			// Check if user has admin permissions (can manage roles)
			hasPermission, err := m.authService.HasPermission(roleID, "roles", "create")
			if err != nil {
				return response.InternalServerError(c, "Failed to Check Admin Permissions", err)
			}

			if !hasPermission {
				return response.Error(c, http.StatusForbidden, "Admin Access Required", nil)
			}

			return next(c)
		}
	}
}
