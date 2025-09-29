package controller

import (
	"log"
	"time"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
	"github.com/DieGopherLT/LatensBackend/internal/services/github"
	"github.com/DieGopherLT/LatensBackend/internal/services/token"
	"github.com/DieGopherLT/LatensBackend/internal/services/users"
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
		log.Println("Error parsing login body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Check if user exists
	user, err := h.UserService.GetUserByGitHubID(c.Context(), body.GithubId)
	if err != nil {
		log.Println("Error fetching user by GitHub ID:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch user",
			"details": err.Error(),
		})
	}

	valid, err := h.GitHubService.ValidateToken(body.AccessToken)
	if err != nil || !valid {
		log.Println("Error validating GitHub token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Could not validate github token",
			"details": err.Error(),
		})
	}

	if user == nil {
		user = &models.UserDocument{
			GithubID:    body.GithubId,
			Username:    body.Username,
			Name:        body.Name,
			Email:       body.Email,
			AccessToken: body.AccessToken,
			AvatarURL:   body.AvatarURL,
		}
		if err := h.UserService.CreateUser(c.Context(), user); err != nil {
			log.Println("Error creating user:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to create user",
				"details": err.Error(),
			})
		}

		payload := token.Payload{
			UserID: user.ID.Hex(),
			Email:  user.Email,
		}

		jwtToken, err := token.Sign(payload)
		if err != nil {
			log.Println("Error signing JWT token:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to generate token",
				"details": err.Error(),
			})
		}
		h.setAuthCookie(c, jwtToken)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
			"user":    user,
		})
	}

	err = h.UserService.UpdateUser(c.Context(), user.ID.Hex(), map[string]any{
		"access_token": body.AccessToken,
	})
	if err != nil {
		log.Println("Error updating user access token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to sync auth",
			"details": err.Error(),
		})
	}

	payload := token.Payload{
		UserID: user.ID.Hex(),
		Email:  user.Email,
	}

	jwtToken, err := token.Sign(payload)
	if err != nil {
		log.Println("Error signing JWT token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate token",
			"details": err.Error(),
		})
	}

	h.setAuthCookie(c, jwtToken)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Auth synced successfully",
		"user":    user,
	})
}

func (h *AuthHandler) setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		MaxAge:   int(time.Hour * 24 * 7 / time.Second), // 7 days
		Path:     "/",
	})
}