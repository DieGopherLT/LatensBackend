package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/DieGopherLT/LatensBackend/pkg/graphql"
)

const (
	GITHUB_PING_ENDPOINT = "https://api.github.com/octocat"
	GITHUB_GRAPHQL_ENDPOINT = "https://api.github.com/graphql"
)

type GithubService struct {
	gqlClient *graphql.Client
}

// NewGithubService creates a new GitHub service without GraphQL client
func NewGithubService() *GithubService {
	return &GithubService{
		gqlClient: nil,
	}
}

// NewGithubServiceWithToken creates a new GitHub service with authenticated GraphQL client
func NewGithubServiceWithToken(token string) *GithubService {
	gqlClient := graphql.NewAuthenticatedClient(GITHUB_GRAPHQL_ENDPOINT, token)

	return &GithubService{
		gqlClient: gqlClient,
	}
}

func (s *GithubService) ValidateToken(token string) (bool, error) {
	req, err := http.NewRequest("GET", GITHUB_PING_ENDPOINT, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode != http.StatusOK {
		return false, fmt.Errorf("GitHub token validation failed with status code: %d", statusCode)
	}

	return true, nil
}

// GetRepositoryMetadata fetches essential repository metadata for database storage
func (s *GithubService) GetRepositoryMetadata(ctx context.Context, token, owner, name string) (*RepositoryMetadataResponse, error) {
	gqlClient := graphql.NewAuthenticatedClient(GITHUB_GRAPHQL_ENDPOINT, token)

	variables := map[string]any{
		"owner": owner,
		"name":  name,
	}

	return graphql.ExecuteQuery[RepositoryMetadataResponse](gqlClient, ctx, RepositoryMetadataQuery, variables)
}

// GetInformationForAwakening fetches data for sleep score calculation
func (s *GithubService) GetInformationForAwakening(ctx context.Context, token, owner, name string, since time.Time) (*SleepAnalysisResponse, error) {
	gqlClient := graphql.NewAuthenticatedClient(GITHUB_GRAPHQL_ENDPOINT, token)

	variables := map[string]any{
		"owner": owner,
		"name":  name,
		"since": since.Format(time.RFC3339),
	}

	return graphql.ExecuteQuery[SleepAnalysisResponse](gqlClient, ctx, SleepAnalysisQuery, variables)
}

// GetUserRepositories fetches only repositories owned by the authenticated user (basic tier)
func (s *GithubService) GetUserRepositories(ctx context.Context, token string, first int, after *string) (*OwnedRepositoriesResponse, error) {
	gqlClient := graphql.NewAuthenticatedClient(GITHUB_GRAPHQL_ENDPOINT, token)

	variables := map[string]any{
		"first": first,
	}

	// Add cursor for pagination if provided
	if after != nil {
		variables["after"] = *after
	}

	return graphql.ExecuteQuery[OwnedRepositoriesResponse](gqlClient, ctx, OwnedRepositoriesQuery, variables)
}
