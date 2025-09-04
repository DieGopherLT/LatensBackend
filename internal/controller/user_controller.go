package controller

import (
	"github.com/DieGopherLT/mfc_backend/internal/database/repository"
	"github.com/gofiber/fiber/v2"
)

// Handlers HTTP Request Handlers
type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	return c.SendString("Create User")
}
