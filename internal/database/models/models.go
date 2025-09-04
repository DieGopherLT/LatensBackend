package models

// User stores personal and github related information about a user
type User struct {
	ID             string `bson:"_id,omitempty" json:"id,omitempty"`
	GitHubUsername string `bson:"github_username" json:"github_username"`
	Email          string `bson:"email" json:"email"`
	FullName       string `bson:"full_name" json:"full_name"`
	AvatarURL      string `bson:"avatar_url" json:"avatar_url"`
	AccessToken    string `bson:"access_token" json:"access_token"`
}
