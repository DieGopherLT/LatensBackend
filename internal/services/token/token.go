package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID            string `json:"user_id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	GitHubAccessToken string `json:"github_access_token"`
}

func Sign(payload Payload) (string, error) {
	claims := jwt.MapClaims{
		"user_id":             payload.UserID,
		"username":            payload.Username,
		"email":               payload.Email,
		"github_access_token": payload.GitHubAccessToken,
		"exp":                 time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":                 time.Now().Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	return jwtToken.SignedString([]byte(jwtSecret))
}

func Parse(tokenString string) (*Payload, error) {
	callback := func(t *jwt.Token) (any, error) {
		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	}

	parsedToken, err := jwt.Parse(tokenString, callback, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	userID, _ := claims["user_id"].(string)
	username, _ := claims["username"].(string)
	email, _ := claims["email"].(string)
	githubToken, _ := claims["github_access_token"].(string)

	payload := &Payload{
		UserID:            userID,
		Username:          username,
		Email:             email,
		GitHubAccessToken: githubToken,
	}

	return payload, nil
}
