package controller

import (
	"log"
	"time"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
	"github.com/DieGopherLT/LatensBackend/internal/services/github"
	"github.com/DieGopherLT/LatensBackend/internal/services/repos"
	"github.com/DieGopherLT/LatensBackend/internal/services/sleep"
	"github.com/DieGopherLT/LatensBackend/internal/services/token"
	"github.com/DieGopherLT/LatensBackend/internal/services/users"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ReposHandler handles requests related to GitHub repositories
type ReposHandler struct {
	reposService  *repos.ReposService
	githubService *github.GithubService
	userService   *users.UserService
}

func NewReposHandler(reposService *repos.ReposService, githubService *github.GithubService, userService *users.UserService) *ReposHandler {
	return &ReposHandler{reposService: reposService, githubService: githubService, userService: userService}
}

func (h *ReposHandler) GetRepos(c *fiber.Ctx) error {
	user := c.Locals("user").(token.Payload)

	repos, err := h.reposService.GetRepositoriesByUserID(c.Context(), user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch repositories. Please try later.",
		})
	}

	if len(repos) == 0 {
		repos = []*models.RepositoryDocument{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repos": repos,
	})
}

func (h *ReposHandler) SyncRepos(c *fiber.Ctx) error{
	user := c.Locals("user").(token.Payload)

	var repos []*models.RepositoryDocument
	var after string
	first := 25

	userGithubToken, err := h.userService.GetUserGitHubToken(c.Context(), user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user GitHub token. Please try later.",
		})
	}

	userID, err := bson.ObjectIDFromHex(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user ID. Please try later.",
			"details": err.Error(),
		})
	}

	for {
		response, err := h.githubService.GetUserRepositories(c.Context(), userGithubToken, first, after)
		if err != nil {
			log.Println("Error fetching user repositories from GitHub:", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to sync repositories from GitHub. Please try later",
				"details": err.Error(),
			})
		}

		log.Println("Fetched", len(response.Viewer.Repositories.Nodes), "repositories from GitHub for user", user.UserID)

		after = response.Viewer.Repositories.PageInfo.EndCursor
		syncTime := time.Now()
		newRepos := lo.Map(response.Viewer.Repositories.Nodes, func(repo github.OwnedRepository, _ int) *models.RepositoryDocument {
			return &models.RepositoryDocument{
				GitHubID:    repo.ID,
				UserID:      userID,
				SleepScore:  sleep.CalculateScore(&repo),
				Name:        repo.Name,
				FullName:    repo.NameWithOwner,
				Description: repo.Description,
				IsPrivate:   repo.IsPrivate,
				IsFork:      repo.IsFork,
				IsDisabled:  repo.IsDisabled,
				IsArchived:  repo.IsArchived,
				URL:         repo.URL,
				DefaultBranch: models.DefaultBranch{
					Name:          repo.DefaultBranchRef.Name,
					CommittedDate: repo.DefaultBranchRef.Target.CommittedDate.String(),
					Author:        repo.DefaultBranchRef.Target.Author.Name,
				},
				CreatedAt:         repo.CreatedAt.String(),
				UpdatedAt:         repo.UpdatedAt.String(),
				PushedAt:          repo.PushedAt.String(),
				IssuesCount:       repo.Issues.TotalCount,
				PullRequestsCount: repo.PullRequests.TotalCount,
				PrimaryLanguage: models.PrimaryLanguage{
					Name:  repo.PrimaryLanguage.Name,
					Color: repo.PrimaryLanguage.Color,
				},
				License: repo.LicenseInfo.Name,
				Metadata: models.RepositoryMetadata{
					SyncedAt: syncTime,
				},
			}
		})

		log.Println("Fetched", len(newRepos), "repositories from GitHub for user", user.UserID)

		repos = append(repos, newRepos...)

		if !response.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}
	}

	log.Println("Synced", len(repos), "repositories for user", user.UserID)

	if len(repos) > 0 {
		err := h.reposService.CreateManyRepositories(c.Context(), repos)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to sync repositories. Please try later.",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Repositories synced successfully.",
		"count":   len(repos),
		"repos":   repos,
	})	
}
