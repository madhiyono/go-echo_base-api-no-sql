package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/pkg/response"
	"github.com/madhiyono/base-api-nosql/pkg/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRole creates a new role (admin only)
func (h *RoleHandler) CreateRole(c echo.Context) error {
	role := new(models.Role)
	if err := c.Bind(role); err != nil {
		h.logger.Error("Failed to Bind Role: %v", err)
		return response.BadRequest(c, "Failed to Create Role: Invalid Request Format", nil)
	}

	// Validate role data
	if err := validation.ValidateStruct(role); err != nil {
		validationErrors := validation.ValidateStructDetailed(role)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation Error for Role: %s", vErr)
		}
		return response.BadRequest(c, "Failed to Create Role: Validation Error", nil)
	}

	if err := h.roleRepo.Create(role); err != nil {
		h.logger.Error("Failed to Create Role: %v", err)
		return response.InternalServerError(c, "Failed to Create Role", err)
	}

	return response.Created(c, "Role Created Successfully", role)
}

// GetRole retrieves a role by ID (admin only)
func (h *RoleHandler) GetRole(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid Role ID", nil)
	}

	role, err := h.roleRepo.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to Get Role: %v", err)
		return response.NotFound(c, "Role Not Found")
	}

	return response.Success(c, "Role Retrieved Successfully", role)
}

// UpdateRole updates an existing role (admin only)
func (h *RoleHandler) UpdateRole(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid Role ID", nil)
	}

	role := new(models.Role)
	if err := c.Bind(role); err != nil {
		h.logger.Error("Failed to Bind Role: %v", err)
		return response.BadRequest(c, "Failed to Update Role: Invalid Request Format", nil)
	}

	// Validate role data
	if err := validation.ValidateStruct(role); err != nil {
		validationErrors := validation.ValidateStructDetailed(role)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation Error for Role: %s", vErr)
		}
		return response.BadRequest(c, "Failed to Update Role: Validation Error", nil)
	}

	if err := h.roleRepo.Update(id, role); err != nil {
		h.logger.Error("Failed to Update Role: %v", err)
		return response.InternalServerError(c, "Failed to Update Role", err)
	}

	return response.Success(c, "Role Updated Successfully", role)
}

// DeleteRole deletes a role by ID (admin only)
func (h *RoleHandler) DeleteRole(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid Role ID", nil)
	}

	if err := h.roleRepo.Delete(id); err != nil {
		h.logger.Error("Failed to Delete Role: %v", err)
		return response.InternalServerError(c, "Failed to Delete Role", err)
	}

	return response.Success(c, "Role Deleted Successfully", nil)
}

// ListRoles returns all roles (admin only)
func (h *RoleHandler) ListRoles(c echo.Context) error {
	roles, err := h.roleRepo.List()
	if err != nil {
		h.logger.Error("Failed to List Roles: %v", err)
		return response.InternalServerError(c, "Failed to Retrieve Roles", err)
	}

	return response.Success(c, "Roles Retrieved Successfully", roles)
}
