package models

// User stores personal and github related information about a user
type User struct {
	ID          string `bson:"_id,omitempty" json:"id,omitempty"`
	GithubID    string `bson:"github_id" json:"github_id"`
	Username    string `bson:"username" json:"username"`
	Name        string `bson:"name" json:"name"`
	Email       string `bson:"email" json:"email"`
	AvatarURL   string `bson:"avatar_url" json:"avatar_url"`
	AccessToken string `bson:"access_token" json:"access_token"`
}
