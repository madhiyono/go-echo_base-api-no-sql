package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/pkg/response"
	"github.com/madhiyono/base-api-nosql/pkg/validation"
)

// Register creates a new user account
func (h *AuthHandler) Register(c echo.Context) error {
	request := new(models.RegisterRequest)
	if err := c.Bind(request); err != nil {
		h.logger.Error("Failed to bind register request: %v", err)
		return response.BadRequest(c, "Failed to Register: Invalid Request Format", nil)
	}

	// Validate request
	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.ValidateStructDetailed(request)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation Error for Register: %s", vErr)
		}
		return response.BadRequest(c, "Failed to Register: Validation Error", nil)
	}

	authResponse, err := h.authService.Register(request)
	if err != nil {
		h.logger.Error("Failed to Register User: %v", err)
		return response.BadRequest(c, "Failed to Register: "+err.Error(), nil)
	}

	return response.Created(c, "User Registered Successfully", authResponse)
}

// Login authenticates a user
func (h *AuthHandler) Login(c echo.Context) error {
	request := new(models.LoginRequest)
	if err := c.Bind(request); err != nil {
		h.logger.Error("Failed to Bind Login Request: %v", err)
		return response.BadRequest(c, "Failed to Login: Invalid Request Format", nil)
	}

	// Validate request
	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.ValidateStructDetailed(request)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation Error for Login: %s", vErr)
		}
		return response.BadRequest(c, "Failed to Login: Validation Error", nil)
	}

	authResponse, err := h.authService.Login(request)
	if err != nil {
		h.logger.Error("Failed to Login User: %v", err)
		return response.Error(c, http.StatusUnauthorized, "Failed to Login: Invalid Credentials", nil)
	}

	return response.Success(c, "Login Successful", authResponse)
}

// Register auth routes
func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	authGroup := e.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
	}
}
