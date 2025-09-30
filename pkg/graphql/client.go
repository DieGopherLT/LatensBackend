package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	gqlclient "github.com/hasura/go-graphql-client"
)

// Query represents a GraphQL query string with type safety
type Query string

// Client is a simplified GraphQL client wrapper
type Client struct {
	gqlClient *gqlclient.Client
}

// NewClient creates a new GraphQL client with the given endpoint and HTTP client
func NewClient(endpoint string, httpClient *http.Client) *Client {
	return &Client{
		gqlClient: gqlclient.NewClient(endpoint, httpClient),
	}
}

// ExecuteQuery executes a GraphQL query with variables and returns a typed response
func ExecuteQuery[T any](c *Client, ctx context.Context, query Query, variables map[string]any) (*T, error) {
	if c.gqlClient == nil {
		return nil, fmt.Errorf("GraphQL client not initialized")
	}

	// Execute the query using ExecRaw for string queries
	raw, err := c.gqlClient.ExecRaw(ctx, string(query), variables)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GraphQL query: %w", err)
	}

	// Parse the raw response into the specified type
	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	return &result, nil
}

// ExecuteQueryRaw executes a GraphQL query with variables and returns the raw response
func (c *Client) ExecuteQueryRaw(ctx context.Context, query Query, variables map[string]any) (map[string]any, error) {
	if c.gqlClient == nil {
		return nil, fmt.Errorf("GraphQL client not initialized")
	}

	// Execute the query using ExecRaw for string queries
	raw, err := c.gqlClient.ExecRaw(ctx, string(query), variables)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GraphQL query: %w", err)
	}

	// Parse the raw response
	var result map[string]any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	return result, nil
}


// ExecuteMutation executes a GraphQL mutation with variables and returns a typed response
func ExecuteMutation[T any](c *Client, ctx context.Context, mutation Query, variables map[string]any) (*T, error) {
	if c.gqlClient == nil {
		return nil, fmt.Errorf("GraphQL client not initialized")
	}

	// Execute the mutation using ExecRaw for string mutations
	raw, err := c.gqlClient.ExecRaw(ctx, string(mutation), variables)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GraphQL mutation: %w", err)
	}

	// Parse the raw response into the specified type
	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	return &result, nil
}

// ExecuteMutationRaw executes a GraphQL mutation with variables and returns the raw response
func (c *Client) ExecuteMutationRaw(ctx context.Context, mutation Query, variables map[string]any) (map[string]any, error) {
	if c.gqlClient == nil {
		return nil, fmt.Errorf("GraphQL client not initialized")
	}

	// Execute the mutation using ExecRaw for string mutations
	raw, err := c.gqlClient.ExecRaw(ctx, string(mutation), variables)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GraphQL mutation: %w", err)
	}

	// Parse the raw response
	var result map[string]any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	return result, nil
}

// AuthedTransport implements http.RoundTripper to add authentication headers
type AuthedTransport struct {
	Token   string
	Wrapped http.RoundTripper
}

// RoundTrip adds the Authorization header to each request
func (t *AuthedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))

	wrapped := t.Wrapped
	if wrapped == nil {
		wrapped = http.DefaultTransport
	}

	return wrapped.RoundTrip(req)
}

// NewAuthenticatedClient creates a new GraphQL client with token authentication
func NewAuthenticatedClient(endpoint, token string) *Client {
	httpClient := &http.Client{
		Transport: &AuthedTransport{
			Token: token,
		},
	}

	return NewClient(endpoint, httpClient)
}