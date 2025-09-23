package repository

import (
	"context"

	"github.com/DieGopherLT/mfc_backend/internal/database/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByGitHubID(ctx context.Context, githubID string) (*models.User, error)
	Update(ctx context.Context, id string, update map[string]any) error
}

type GitHubReposRepository interface {
	Create(ctx context.Context, repo *models.GitHubRepository) error
	CreateMany(ctx context.Context, repos []*models.GitHubRepository) error
}