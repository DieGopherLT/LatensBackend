package repository

import (
	"context"
	"errors"

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

func (r *MongoUserRepository) FindByGitHubID(ctx context.Context, githubID string) (*models.User, error) {
	var user models.User

	condition := map[string]any{"github_id": githubID}
	err := r.collection.FindOne(ctx, condition).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, id string, update map[string]any) error {
	condition := map[string]any{"_id": id}
	_, err := r.collection.UpdateOne(ctx, condition, map[string]any{"$set": update})
	return err
}
