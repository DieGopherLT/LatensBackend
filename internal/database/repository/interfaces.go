package repository

import (
	"context"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByGitHubID(ctx context.Context, githubID string) (*models.User, error)
	FindAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, id string, update map[string]any) error
	Delete(ctx context.Context, id string) error
}

type GitHubReposRepository interface {
	Create(ctx context.Context, repo *models.GitHubRepository) error
	CreateMany(ctx context.Context, repos []*models.GitHubRepository) error
	FindByID(ctx context.Context, id string, userID string) (*models.GitHubRepository, error)
	FindAllByUser(ctx context.Context, userID string) ([]*models.GitHubRepository, error)
	Update(ctx context.Context, id string, userID string, update map[string]any) error
	Delete(ctx context.Context, id string, userID string) error
}