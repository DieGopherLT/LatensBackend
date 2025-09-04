package repository

import (
	"context"

	"github.com/DieGopherLT/mfc_backend/internal/database/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
}
