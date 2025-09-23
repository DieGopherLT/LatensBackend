package repository

import (
	"context"

	"github.com/DieGopherLT/mfc_backend/internal/database/models"
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
	docs := make([]models.GitHubRepository, len(repos))
	for i, repo := range repos {
		docs[i] = *repo
	}
	_, err := r.collection.InsertMany(ctx, docs)
	return err
}
