package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission struct {
	Resource string `json:"resource" bson:"resource"`
	Action   string `json:"action" bson:"action"`
}

type Role struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Description string             `json:"description" bson:"description"`
	Permissions []Permission       `json:"permissions" bson:"permissions"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Predefined permissions
const (
	ResourceUsers = "users"
	ResourceRoles = "roles"

	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// Helper function to create permission
func NewPermission(resource, action string) Permission {
	return Permission{
		Resource: resource,
		Action:   action,
	}
}
