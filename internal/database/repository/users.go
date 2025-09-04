package repository

import (
	"context"

	"github.com/DieGopherLT/mfc_backend/internal/database/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Database) *MongoUserRepository {
	collection := client.Collection("users")
	return &MongoUserRepository{collection: collection}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}
