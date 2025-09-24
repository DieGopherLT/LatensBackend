package repository

import (
	"context"
	"errors"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoGitHubReposRepository struct {
	collection *mongo.Collection
}

func NewGitHubReposRepository(client *mongo.Database) *MongoGitHubReposRepository {
	collection := client.Collection("github_repositories")
	return &MongoGitHubReposRepository{collection: collection}
}

func (r *MongoGitHubReposRepository) Create(ctx context.Context, repo *models.GitHubRepository) error {
	_, err := r.collection.InsertOne(ctx, repo)
	return err
}

func (r *MongoGitHubReposRepository) CreateMany(ctx context.Context, repos []*models.GitHubRepository) error {
	docs := make([]any, len(repos))
	for i, repo := range repos {
		docs[i] = repo
	}
	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

func (r *MongoGitHubReposRepository) FindByID(ctx context.Context, id string) (*models.GitHubRepository, error) {
	var repo models.GitHubRepository

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	condition := map[string]any{"_id": objectID}
	err = r.collection.FindOne(ctx, condition).Decode(&repo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *MongoGitHubReposRepository) FindByUserID(ctx context.Context, userID string) ([]*models.GitHubRepository, error) {
	condition := map[string]any{"user_id": userID}
	cursor, err := r.collection.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repos []*models.GitHubRepository
	for cursor.Next(ctx) {
		var repo models.GitHubRepository
		if err := cursor.Decode(&repo); err != nil {
			return nil, err
		}
		repos = append(repos, &repo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *MongoGitHubReposRepository) FindAll(ctx context.Context) ([]*models.GitHubRepository, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repos []*models.GitHubRepository
	for cursor.Next(ctx) {
		var repo models.GitHubRepository
		if err := cursor.Decode(&repo); err != nil {
			return nil, err
		}
		repos = append(repos, &repo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *MongoGitHubReposRepository) Update(ctx context.Context, id string, update map[string]any) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	condition := map[string]any{"_id": objectID}
	_, err = r.collection.UpdateOne(ctx, condition, map[string]any{"$set": update})
	return err
}

func (r *MongoGitHubReposRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	condition := map[string]any{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, condition)
	return err
}
