package mongo

import (
	"context"
	"time"

	"github.com/madhiyono/base-api-nosql/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type roleRepository struct {
	collection *mongo.Collection
}

func NewRoleRepository(db *mongo.Database) *roleRepository {
	return &roleRepository{
		collection: db.Collection("roles"),
	}
}

func (r *roleRepository) Create(role *models.Role) error {
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(context.TODO(), role)
	if err != nil {
		return err
	}

	role.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *roleRepository) GetByID(id primitive.ObjectID) (*models.Role, error) {
	var role models.Role
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) Update(id primitive.ObjectID, role *models.Role) error {
	role.UpdatedAt = time.Now()
	role.ID = id

	filter := bson.M{"_id": id}
	update := bson.M{"$set": role}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *roleRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (r *roleRepository) List() ([]*models.Role, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var roles []*models.Role
	for cursor.Next(context.TODO()) {
		var role models.Role
		if err := cursor.Decode(&role); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *roleRepository) HasPermission(roleID primitive.ObjectID, resource, action string) (bool, error) {
	var role models.Role
	err := r.collection.FindOne(context.TODO(), bson.M{
		"_id":       roleID,
		"is_active": true,
	}).Decode(&role)

	if err != nil {
		return false, err
	}

	permission := models.NewPermission(resource, action)
	for _, p := range role.Permissions {
		if p.Resource == permission.Resource && p.Action == permission.Action {
			return true, nil
		}
	}

	return false, nil
}
