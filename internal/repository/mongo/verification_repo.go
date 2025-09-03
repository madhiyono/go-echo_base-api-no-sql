package mongo

import (
	"context"
	"time"

	"github.com/madhiyono/base-api-nosql/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type verificationRepository struct {
	collection *mongo.Collection
}

func NewVerificationRepository(db *mongo.Database) *verificationRepository {
	return &verificationRepository{
		collection: db.Collection("email_verifications"),
	}
}

func (r *verificationRepository) Create(verification *models.EmailVerification) error {
	verification.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(context.TODO(), verification)
	if err != nil {
		return err
	}

	verification.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *verificationRepository) GetByToken(token string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.collection.FindOne(context.TODO(), bson.M{
		"token":      token,
		"is_used":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&verification)

	if err != nil {
		return nil, err
	}

	return &verification, nil
}

func (r *verificationRepository) MarkAsUsed(id primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"is_used": true,
			"used_at": &now,
		},
	}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *verificationRepository) GetByUserID(userID primitive.ObjectID) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.collection.FindOne(context.TODO(), bson.M{
		"user_id":    userID,
		"is_used":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&verification)

	if err != nil {
		return nil, err
	}

	return &verification, nil
}
