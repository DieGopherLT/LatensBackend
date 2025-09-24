package controller

import (
	"github.com/DieGopherLT/mfc_backend/internal/database/models"
	"github.com/DieGopherLT/mfc_backend/internal/services/github"
	"github.com/DieGopherLT/mfc_backend/internal/services/token"
	"github.com/DieGopherLT/mfc_backend/internal/services/users"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	UserService   *users.UserService
	GitHubService *github.GithubService
}

func NewAuthHandler(userService *users.UserService, githubService *github.GithubService) *AuthHandler {
	return &AuthHandler{
		UserService:   userService,
		GitHubService: githubService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		GithubId    string `json:"github_id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		AccessToken string `json:"access_token"`
		AvatarURL   string `json:"avatar_url"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Check if user exists
	user, err := h.UserService.GetUserByGitHubID(c.Context(), body.GithubId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch user",
			"details": err.Error(),
		})
	}

	valid, err := h.GitHubService.ValidateToken(body.AccessToken)
	if err != nil || !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Could not validate github token",
			"details": err.Error(),
		})
	}

	if user == nil {
		user = &models.User{
			GithubID:    body.GithubId,
			Username:    body.Username,
			Name:        body.Name,
			Email:       body.Email,
			AccessToken: body.AccessToken,
			AvatarURL:   body.AvatarURL,
		}
		if err := h.UserService.CreateUser(c.Context(), user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to create user",
				"details": err.Error(),
			})
		}

		payload := token.Payload{
			UserID:            user.ID,
			Username:          user.Username,
			Email:             user.Email,
			GitHubAccessToken: user.AccessToken,
		}

		jwtToken, err := token.Sign(payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to generate token",
				"details": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
			"user":    user,
			"token":   jwtToken,
		})
	}

	err = h.UserService.UpdateUser(c.Context(), user.ID, map[string]any{
		"access_token": body.AccessToken,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to sync auth",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Auth synced successfully",
		"user":    user,
	})
}
