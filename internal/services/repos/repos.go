package repos

import (
	"context"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
	"github.com/DieGopherLT/LatensBackend/internal/database/repository"
)

type ReposService struct {
	repo repository.GitHubReposRepository
}

func NewReposService(repo repository.GitHubReposRepository) *ReposService {
	return &ReposService{repo: repo}
}

func (s *ReposService) CreateRepository(ctx context.Context, repo *models.RepositoryDocument) error {
	return s.repo.Create(ctx, repo)
}

func (s *ReposService) CreateManyRepositories(ctx context.Context, repos []*models.RepositoryDocument) error {
	return s.repo.CreateMany(ctx, repos)
}

func (s *ReposService) GetRepositoryByID(ctx context.Context, id string, userID string) (*models.RepositoryDocument, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *ReposService) GetRepositoriesByUserID(ctx context.Context, userID string) ([]*models.RepositoryDocument, error) {
	return s.repo.FindAllByUser(ctx, userID)
}

func (s *ReposService) UpdateRepository(ctx context.Context, id string, userID string, update map[string]any) error {
	return s.repo.Update(ctx, id, userID, update)
}

func (s *ReposService) DeleteRepository(ctx context.Context, id string, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}