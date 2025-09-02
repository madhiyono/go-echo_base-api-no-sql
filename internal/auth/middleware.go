package auth

import (
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/pkg/response"
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
		c.Set("role", claims.Role)

		return next(c)
	}
}

// RequireRole middleware for role-based authorization
func (m *Middleware) RequireRole(roles ...models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role").(models.UserRole)

			// Check if user has required role
			if slices.Contains(roles, userRole) {
				return next(c)
			}

			return response.Error(c, http.StatusForbidden, "Insufficient Permissions", nil)
		}
	}
}

// RequireAdmin middleware for admin-only access
func (m *Middleware) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return m.RequireRole(models.RoleAdmin)(next)
}

// RequireUser middleware for user access
func (m *Middleware) RequireUser(next echo.HandlerFunc) echo.HandlerFunc {
	return m.RequireRole(models.RoleAdmin, models.RoleUser)(next)
}
