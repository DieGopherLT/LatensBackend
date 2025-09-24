package repos

import (
	"context"

	"github.com/DieGopherLT/mfc_backend/internal/database/models"
	"github.com/DieGopherLT/mfc_backend/internal/database/repository"
)

type ReposService struct {
	repo repository.GitHubReposRepository
}

func NewReposService(repo repository.GitHubReposRepository) *ReposService {
	return &ReposService{repo: repo}
}

func (s *ReposService) CreateRepository(ctx context.Context, repo *models.GitHubRepository) error {
	return s.repo.Create(ctx, repo)
}

func (s *ReposService) CreateManyRepositories(ctx context.Context, repos []*models.GitHubRepository) error {
	return s.repo.CreateMany(ctx, repos)
}

func (s *ReposService) GetRepositoryByID(ctx context.Context, id string) (*models.GitHubRepository, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ReposService) GetRepositoriesByUserID(ctx context.Context, userID string) ([]*models.GitHubRepository, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *ReposService) GetAllRepositories(ctx context.Context) ([]*models.GitHubRepository, error) {
	return s.repo.FindAll(ctx)
}

func (s *ReposService) UpdateRepository(ctx context.Context, id string, update map[string]any) error {
	return s.repo.Update(ctx, id, update)
}

func (s *ReposService) DeleteRepository(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}