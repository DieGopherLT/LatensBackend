package controller

import (
	"github.com/DieGopherLT/LatensBackend/internal/database/models"
	"github.com/DieGopherLT/LatensBackend/internal/services/github"
	"github.com/DieGopherLT/LatensBackend/internal/services/repos"
	"github.com/DieGopherLT/LatensBackend/internal/services/token"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// ReposHandler handles requests related to GitHub repositories
type ReposHandler struct {
	reposService  *repos.ReposService
	githubService *github.GithubService
}

func NewReposHandler(reposService *repos.ReposService, githubService *github.GithubService) *ReposHandler {
	return &ReposHandler{reposService: reposService, githubService: githubService}
}

func (h *ReposHandler) GetRepos(c *fiber.Ctx) error {
	repos, err := h.reposService.GetAllRepositories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch repositories. Please try later.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repos": repos,
	})
}

func (h *ReposHandler) SyncRepos(c *fiber.Ctx) error{
	user := c.Locals("user").(token.Payload)	
	
	var repos []*models.GitHubRepository
	var after *string
	first := 25

	for {
		response, err := h.githubService.GetUserRepositories(c.Context(), user.GitHubAccessToken, first, after)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to sync repositories from GitHub. Please try later",
				"details": err.Error(),
			})
		}

		after = &response.Data.Viewer.Repositories.PageInfo.EndCursor
		newRepos := lo.Map(response.Data.Viewer.Repositories.Nodes, func(repo github.OwnedRepository, _ int) *models.GitHubRepository {
			return &models.GitHubRepository{
				GitHubID:    repo.ID,
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
			}
		})
		repos = append(repos, newRepos...)

		if !response.Data.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}
	}

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
