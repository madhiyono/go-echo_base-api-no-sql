package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/models"
)

// CreateUser: Creates a new user
func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		h.logger.Error("Failed to Bind User Data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request!"})
	}

	if err := h.userRepo.Create(user); err != nil {
		h.logger.Error("Failed to Create User: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to Create User!"})
	}

	return c.JSON(http.StatusCreated, user)
}

// GetUser: Retrieves a user by ID
func (h *UserHandler) GetUser(c echo.Context) error {
    id := c.Param("id")

    user, err := h.userRepo.GetByID(id)
    if err != nil {
        h.logger.Error("Failed to Get User: %v", err)
        return c.JSON(http.StatusNotFound, map[string]string{"error": "User Not Found!"})
    }

    return c.JSON(http.StatusOK, user)
}

// UpdateUser: Updates an existing user
func (h *UserHandler) UpdateUser(c echo.Context) error {
    id := c.Param("id")

    user := new(models.User)
    if err := c.Bind(user); err != nil {
        h.logger.Error("Failed to Bind User: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request!"})
    }

    if err := h.userRepo.Update(id, user); err != nil {
        h.logger.Error("Failed to Update User: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to Update User!"})
    }

    return c.JSON(http.StatusOK, user)
}

// DeleteUser: Deletes a user by ID
func (h *UserHandler) DeleteUser(c echo.Context) error {
    id := c.Param("id")

    if err := h.userRepo.Delete(id); err != nil {
        h.logger.Error("Failed to Delete User: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to Delete User!"})
    }

    return c.NoContent(http.StatusNoContent)
}

// ListUsers: Returns all users
func (h *UserHandler) ListUsers(c echo.Context) error {
    users, err := h.userRepo.List()
    if err != nil {
        h.logger.Error("Failed to List Users: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to Retrieve Users!"})
    }

    return c.JSON(http.StatusOK, users)
}