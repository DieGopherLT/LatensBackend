package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/DieGopherLT/mfc_backend/pkg/graphql"
)

const (
	GITHUB_PING_ENDPOINT = "https://api.github.com/octocat"
	GITHUB_GRAPHQL_ENDPOINT = "https://api.github.com/graphql"
)

type GithubService struct {
	gqlClient *graphql.Client
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
func (s *GithubService) GetRepositoryMetadata(ctx context.Context, owner, name string) (*RepositoryMetadataResponse, error) {
	if s.gqlClient == nil {
		return nil, fmt.Errorf("GraphQL client not initialized. Use NewGithubServiceWithToken")
	}

	variables := map[string]any{
		"owner": owner,
		"name":  name,
	}

	return graphql.ExecuteQuery[RepositoryMetadataResponse](s.gqlClient, ctx, RepositoryMetadataQuery, variables)
}

// GetInformationForAwakening fetches data for sleep score calculation
func (s *GithubService) GetInformationForAwakening(ctx context.Context, owner, name string, since time.Time) (*SleepAnalysisResponse, error) {
	if s.gqlClient == nil {
		return nil, fmt.Errorf("GraphQL client not initialized. Use NewGithubServiceWithToken")
	}

	variables := map[string]any{
		"owner": owner,
		"name":  name,
		"since": since.Format(time.RFC3339),
	}

	return graphql.ExecuteQuery[SleepAnalysisResponse](s.gqlClient, ctx, SleepAnalysisQuery, variables)
}

// GetUserRepositories fetches only repositories owned by the authenticated user (basic tier)
func (s *GithubService) GetUserRepositories(ctx context.Context, first int, after *string) (*OwnedRepositoriesResponse, error) {
	if s.gqlClient == nil {
		return nil, fmt.Errorf("GraphQL client not initialized. Use NewGithubServiceWithToken")
	}

	variables := map[string]any{
		"first": first,
	}

	// Add cursor for pagination if provided
	if after != nil {
		variables["after"] = *after
	}

	return graphql.ExecuteQuery[OwnedRepositoriesResponse](s.gqlClient, ctx, OwnedRepositoriesQuery, variables)
}
