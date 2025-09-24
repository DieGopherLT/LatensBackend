package controller

import (
	"github.com/DieGopherLT/mfc_backend/internal/database/models"
	"github.com/DieGopherLT/mfc_backend/internal/services/github"
	"github.com/DieGopherLT/mfc_backend/internal/services/repos"
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

func (h *ReposHandler) SyncRepositories(c *fiber.Ctx) error {
	return nil
}

func (h *ReposHandler) startSyncJob(c *fiber.Ctx) {
	var repos []*models.GitHubRepository
	var after *string
	first := 25

	for {
		response, err := h.githubService.GetUserRepositories(c.Context(), body.AccessToken, first, after)
		if err != nil {
			break
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
		_ = h.reposService.CreateManyRepositories(c.Context(), repos)
	}
}
