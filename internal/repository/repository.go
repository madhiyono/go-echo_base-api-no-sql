package repository

import (
	"github.com/madhiyono/base-api-nosql/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	Update(id string, user *models.User) error
	Delete(id string) error
	List() ([]*models.User, error)
}

type AuthRepository interface {
	Create(auth *models.UserAuth) error
	GetByEmail(email string) (*models.UserAuth, error)
	GetByUserID(userID primitive.ObjectID) (*models.UserAuth, error)
	UpdatePassword(userID primitive.ObjectID, password string) error
}
