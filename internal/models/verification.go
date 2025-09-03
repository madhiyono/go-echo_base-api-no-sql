package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailVerification struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Email     string             `json:"email" bson:"email"`
	Token     string             `json:"token" bson:"token"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	IsUsed    bool               `json:"is_used" bson:"is_used"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UsedAt    *time.Time         `json:"used_at,omitempty" bson:"used_at,omitempty"`
}
