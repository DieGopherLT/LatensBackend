package github

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const GITHUB_PING_ENDPOINT = "https://api.github.com/octocat"

type GithubService struct {
}

func NewGithubService() *GithubService {
	return &GithubService{}
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

func (s *GithubService) SyncRepositories(ctx context.Context, token, userID string) {}
