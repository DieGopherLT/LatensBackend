package repository

import (
	"context"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.UserDocument) error
	FindByID(ctx context.Context, id string) (*models.UserDocument, error)
	FindByGitHubID(ctx context.Context, githubID string) (*models.UserDocument, error)
	FindAll(ctx context.Context) ([]*models.UserDocument, error)
	Update(ctx context.Context, id string, update map[string]any) error
	Delete(ctx context.Context, id string) error
}

type GitHubReposRepository interface {
	Create(ctx context.Context, repo *models.RepositoryDocument) error
	CreateMany(ctx context.Context, repos []*models.RepositoryDocument) error
	FindByID(ctx context.Context, id string, userID string) (*models.RepositoryDocument, error)
	FindAllByUser(ctx context.Context, userID string) ([]*models.RepositoryDocument, error)
	Update(ctx context.Context, id string, userID string, update map[string]any) error
	Delete(ctx context.Context, id string, userID string) error
}