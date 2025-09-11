package main

import (
	"github.com/DieGopherLT/mfc_backend/internal/controller"
	"github.com/DieGopherLT/mfc_backend/internal/database/repository"
	"github.com/DieGopherLT/mfc_backend/internal/services/github"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupRoutes(app *fiber.App, db *mongo.Database) *fiber.App {

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Services
	githubService := github.NewGithubService()

	// Handlers
	userHandler := controller.NewUserHandler(userRepo)
	authHandler := controller.NewAuthHandler(userRepo, githubService)

	// Versions
	v1 := app.Group("/api/v1")

	// Groups
	users := v1.Group("/users")
	auth := v1.Group("/auth")

	// User routes
	users.Post("/", userHandler.CreateUser)

	// Auth routes
	auth.Post("/sync", authHandler.Sync)

	return app
}
