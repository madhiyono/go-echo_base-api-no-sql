package mongo

import (
	"context"
	"time"

	"github.com/madhiyono/base-api-nosql/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type authRepository struct {
	collection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) *authRepository {
	return &authRepository{
		collection: db.Collection("user_auth"),
	}
}

func (r *authRepository) Create(auth *models.UserAuth) error {
	auth.CreatedAt = time.Now()
	auth.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(context.TODO(), auth)
	if err != nil {
		return err
	}

	auth.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *authRepository) GetByEmail(email string) (*models.UserAuth, error) {
	var auth models.UserAuth
	err := r.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *authRepository) GetByUserID(userID primitive.ObjectID) (*models.UserAuth, error) {
	var auth models.UserAuth
	err := r.collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *authRepository) UpdatePassword(userID primitive.ObjectID, password string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"password":   password,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *authRepository) UpdateRole(userID, roleID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"role_id":    roleID,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}
