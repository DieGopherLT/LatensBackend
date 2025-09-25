package config

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `validate:"required"`
	MongoDbURI  string `validate:"required"`
	MongoDbName string `validate:"required"`
	OpenAIKey   string `validate:"required"`
	JWTSecret   string `validate:"required"`
}

var (
	ErrEnvLoad = errors.New("error loading .env file")
)

var (
	MongoDBKey  = "MONGODB_URI"
	MongoDbName = "MONGODB_NAME"
	OpenAIKey   = "OPENAI_API_KEY"
	PortKey     = "PORT"
	JWTSecret   = "JWT_SECRET"
)

func New() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, ErrEnvLoad
	}

	config := &Config{
		MongoDbURI:  os.Getenv(MongoDBKey),
		MongoDbName: os.Getenv(MongoDbName),
		OpenAIKey:   os.Getenv(OpenAIKey),
		Port:        os.Getenv(PortKey),
		JWTSecret:   os.Getenv(JWTSecret),
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
