package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/pkg/response"
	"github.com/madhiyono/base-api-nosql/pkg/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// In CreateUser method, you might want to check if the authenticated user has permission
// to create other users (admin-only functionality)
func (h *UserHandler) CreateUser(c echo.Context) error {
	// Check if user has admin role to create other users
	role := c.Get("role").(models.UserRole)
	if role != models.RoleAdmin {
		return response.Error(c, http.StatusForbidden, "Insufficient Permissions to Create Users", nil)
	}

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		h.logger.Error("Failed to Bind User: %v", err)
		return response.BadRequest(c, "Failed to Create User: Invalid Request Format", nil)
	}

	// Validate user data (log detailed errors but return generic message)
	if err := validation.ValidateStruct(user); err != nil {
		// Log detailed validation errors for debugging
		validationErrors := validation.ValidateStructDetailed(user)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation Error for User: %s", vErr)
		}
		return response.BadRequest(c, "Failed to Create User: Validation Error", nil)
	}

	if err := h.userRepo.Create(user); err != nil {
		h.logger.Error("Failed to Create User: %v", err)
		return response.InternalServerError(c, "Failed to Create User: Internal Server Error", nil)
	}

	return response.Created(c, "User Created Successfully", user)
}

// GetUser: Retrieves a user by ID
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to Get User: %v", err)
		return response.NotFound(c, "User Not Found!")
	}

	return response.Success(c, "User Retrieved Successfully", user)
}

// In other methods, you might want to add authorization checks
// For example, users can only update their own profile
func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID := c.Param("id")
	authUserID := c.Get("user_id").(primitive.ObjectID)

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		h.logger.Error("Failed to Bind User: %v", err)
		return response.BadRequest(c, "Failed to Update User: Invalid Request Format", nil)
	}

	// Validate user data (log detailed errors but return generic message)
	if err := validation.ValidateStruct(user); err != nil {
		// Log detailed validation errors for debugging
		validationErrors := validation.ValidateStructDetailed(user)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation Error for User: %s", vErr)
		}
		return response.BadRequest(c, "Failed to Update User: Validation Error", nil)
	}

	// Non-admin users can only update their own profile
	role := c.Get("role").(models.UserRole)
	if role != models.RoleAdmin {
		if userID != authUserID.Hex() {
			return response.Error(c, http.StatusForbidden, "Cannot Update Other Users", nil)
		}
	}

	if err := h.userRepo.Update(userID, user); err != nil {
		h.logger.Error("Failed to Update User: %v", err)
		return response.InternalServerError(c, "Failed to Update User: Internal Server Error", nil)
	}

	return response.Success(c, "User Updated Successfully", user)
}

// DeleteUser: Deletes a user by ID
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	if err := h.userRepo.Delete(id); err != nil {
		h.logger.Error("Failed to Delete User: %v", err)
		return response.InternalServerError(c, "Failed to Delete User: Internal Server Error", nil)
	}

	return response.Success(c, "User Deleted Successfully", nil)
}

// ListUsers: Returns all users
func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userRepo.List()
	if err != nil {
		h.logger.Error("Failed to list users: %v", err)
		return response.InternalServerError(c, "Failed to Retrieve Users: Internal Server Error", nil)
	}

	return response.Success(c, "Users Retrieved Successfully", users)
}
