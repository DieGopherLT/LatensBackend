package models

import "go.mongodb.org/mongo-driver/v2/bson"

// User stores personal and github related information about a user
type User struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	GithubID    string `bson:"github_id" json:"github_id"`
	Username    string `bson:"username" json:"username"`
	Name        string `bson:"name" json:"name"`
	Email       string `bson:"email" json:"email"`
	AvatarURL   string `bson:"avatar_url" json:"avatar_url"`
	AccessToken string `bson:"access_token" json:"access_token"`
}

type GitHubRepository struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	GitHubID          string          `bson:"github_id" json:"github_id"`
	Name              string          `bson:"name" json:"name"`
	FullName          string          `bson:"full_name" json:"full_name"`
	Description       string          `bson:"description" json:"description"`
	IsPrivate         bool            `bson:"is_private" json:"is_private"`
	IsFork            bool            `bson:"is_fork" json:"is_fork"`
	IsDisabled        bool            `bson:"is_disabled" json:"is_disabled"`
	IsArchived        bool            `bson:"is_archived" json:"is_archived"`
	URL               string          `bson:"url" json:"url"`
	DefaultBranch     DefaultBranch   `bson:"default_branch" json:"default_branch"`
	CreatedAt         string          `bson:"created_at" json:"created_at"`
	UpdatedAt         string          `bson:"updated_at" json:"updated_at"`
	PushedAt          string          `bson:"pushed_at" json:"pushed_at"`
	Topics            []string        `bson:"topics" json:"topics"`
	IssuesCount       int             `bson:"issues_count" json:"issues_count"`
	PullRequestsCount int             `bson:"pull_requests_count" json:"pull_requests_count"`
	PrimaryLanguage   PrimaryLanguage `bson:"primary_language" json:"primary_language"`
	License           string          `bson:"license" json:"license"`
}

type DefaultBranch struct {
	Name          string `bson:"name" json:"name"`
	CommittedDate string `bson:"committed_date" json:"committed_date"`
	Author        string `bson:"author" json:"author"`
}

type PrimaryLanguage struct {
	Name  string `bson:"name" json:"name"`
	Color string `bson:"color" json:"color"`
}
