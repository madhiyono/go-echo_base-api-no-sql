package repository

import "github.com/madhiyono/base-api-nosql/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	Update(id string, user *models.User) error
	Delete(id string) error
	List() ([]*models.User, error)
}