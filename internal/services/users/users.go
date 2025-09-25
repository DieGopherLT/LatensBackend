package users

import (
	"context"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
	"github.com/DieGopherLT/LatensBackend/internal/database/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) GetUserByGitHubID(ctx context.Context, githubID string) (*models.User, error) {
	return s.repo.FindByGitHubID(ctx, githubID)
}

func (s *UserService) GetUserGitHubToken(ctx context.Context, id string) (string, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil
	}
	return user.AccessToken, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, update map[string]any) error {
	return s.repo.Update(ctx, id, update)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
