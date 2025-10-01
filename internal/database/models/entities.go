package models

import "go.mongodb.org/mongo-driver/v2/bson"

// UserDocument stores personal and github related information about a user
type UserDocument struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	GithubID    string `bson:"github_id" json:"github_id"`
	Username    string `bson:"username" json:"username"`
	Name        string `bson:"name" json:"name"`
	Email       string `bson:"email" json:"email"`
	AvatarURL   string `bson:"avatar_url" json:"avatar_url"`
	AccessToken string `bson:"access_token" json:"access_token"`
}

type RepositoryDocument struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SleepScore        int             `bson:"sleep_score" json:"sleep_score"`
	UserID            bson.ObjectID `bson:"user_id" json:"user_id"`
	GitHubID          string          `bson:"github_id" json:"github_id"`
	Name              string          `bson:"name" json:"name"`
	FullName          string          `bson:"full_name" json:"full_name"`
	Description       string          `bson:"description" json:"description"`
	IsPrivate         bool            `bson:"is_private" json:"is_private"`
	IsFork            bool            `bson:"is_fork" json:"is_fork"`
	IsDisabled        bool            `bson:"is_disabled" json:"is_disabled"`
	IsArchived        bool            `bson:"is_archived" json:"is_archived"`
	IsHidden          bool            `bson:"is_hidden,omitempty" json:"is_hidden"`
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
	Metadata          RepositoryMetadata `bson:"metadata" json:"metadata"`
}