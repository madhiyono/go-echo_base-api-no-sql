package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/pkg/response"
	"github.com/madhiyono/base-api-nosql/pkg/validation"
)

// Register creates a new user account with email verification
func (h *AuthHandler) Register(c echo.Context) error {
	request := new(models.RegisterRequest)
	if err := c.Bind(request); err != nil {
		h.logger.Error("Failed to bind register request: %v", err)
		return response.BadRequest(c, "Failed to register: Invalid request format", nil)
	}

	// Validate request
	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.ValidateStructDetailed(request)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation error for register: %s", vErr)
		}
		return response.BadRequest(c, "Failed to register: Validation Error", nil)
	}

	// Check if user already exists
	if _, err := h.authRepo.GetByEmail(request.Email); err == nil {
		return response.BadRequest(c, "User already exists", nil)
	}

	// Create user
	user := &models.User{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := h.userRepo.Create(user); err != nil {
		h.logger.Error("Failed to create user: %v", err)
		return response.InternalServerError(c, "Failed to create user", nil)
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(request.Password)
	if err != nil {
		h.logger.Error("Failed to hash password: %v", err)
		return response.InternalServerError(c, "Failed to process registration", nil)
	}

	// Get default role (user role)
	defaultRole, err := h.roleRepo.GetByName("user")
	if err != nil {
		h.logger.Error("Default role not found: %v", err)
		return response.InternalServerError(c, "Failed to process registration", nil)
	}

	// Create auth record
	auth := &models.UserAuth{
		UserID:   user.ID,
		Email:    request.Email,
		Password: hashedPassword,
		RoleID:   defaultRole.ID,
		IsActive: false, // User is inactive until email is verified
	}

	if err := h.authRepo.Create(auth); err != nil {
		h.logger.Error("Failed to create auth record: %v", err)
		return response.InternalServerError(c, "Failed to process registration", nil)
	}

	// Send verification email
	if err := h.emailService.SendVerificationEmail(user.ID, user.Email, user.Name); err != nil {
		h.logger.Error("Failed to send verification email: %v", err)
		// Don't fail registration for email sending error, but log it
	}

	return response.Created(c, "User registered successfully. Please check your email for verification.", nil)
}

// VerifyEmail verifies user email address
func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return response.BadRequest(c, "Verification token is required", nil)
	}

	// Verify email token
	if err := h.emailService.VerifyEmail(token); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid or expired verification token", nil)
	}

	// Get verification record to get user ID
	verification, err := h.verifyRepo.GetByToken(token)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid verification token", nil)
	}

	// Activate user account
	if err := h.authRepo.ActivateUser(verification.UserID); err != nil {
		h.logger.Error("Failed to activate user: %v", err)
		return response.InternalServerError(c, "Failed to activate account", nil)
	}

	return response.Success(c, "Email verified successfully. Your account is now active.", nil)
}

// ResendVerification sends verification email again
func (h *AuthHandler) ResendVerification(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return response.BadRequest(c, "Email is required", nil)
	}

	// Get user auth record
	auth, err := h.authRepo.GetByEmail(email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return response.Success(c, "If the email exists, verification has been sent.", nil)
	}

	// Check if already verified
	if auth.IsActive {
		return response.Success(c, "Email is already verified.", nil)
	}

	// Get user details
	user, err := h.userRepo.GetByID(auth.UserID.Hex())
	if err != nil {
		h.logger.Error("Failed to get user: %v", err)
		return response.InternalServerError(c, "Failed to process request", nil)
	}

	// Send verification email
	if err := h.emailService.SendVerificationEmail(user.ID, user.Email, user.Name); err != nil {
		h.logger.Error("Failed to send verification email: %v", err)
		return response.InternalServerError(c, "Failed to send verification email", nil)
	}

	return response.Success(c, "Verification email sent successfully.", nil)
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
		authGroup.GET("/verify/:token", h.VerifyEmail)
		authGroup.POST("/resend-verification", h.ResendVerification)
		authGroup.POST("/login", h.Login)
	}
}
