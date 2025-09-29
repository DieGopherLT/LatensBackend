package repository

import (
	"context"
	"errors"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Database) *MongoUserRepository {
	collection := client.Collection("users")
	return &MongoUserRepository{collection: collection}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *models.UserDocument) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepository) FindByGitHubID(ctx context.Context, githubID string) (*models.UserDocument, error) {
	var user models.UserDocument

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

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*models.UserDocument, error) {
	var user models.UserDocument

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	condition := map[string]any{"_id": objectID}
	err = r.collection.FindOne(ctx, condition).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) FindAll(ctx context.Context) ([]*models.UserDocument, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.UserDocument
	for cursor.Next(ctx) {
		var user models.UserDocument
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, id string, update map[string]any) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	condition := map[string]any{"_id": objectID}
	_, err = r.collection.UpdateOne(ctx, condition, map[string]any{"$set": update})
	return err
}

func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	condition := map[string]any{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, condition)
	return err
}
