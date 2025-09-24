package main

import (
	"github.com/DieGopherLT/mfc_backend/internal/controller"
	"github.com/DieGopherLT/mfc_backend/internal/database/repository"
	"github.com/DieGopherLT/mfc_backend/internal/services/github"
	"github.com/DieGopherLT/mfc_backend/internal/services/repos"
	"github.com/DieGopherLT/mfc_backend/internal/services/users"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupRoutes(app *fiber.App, db *mongo.Database) *fiber.App {

	// Repositories
	userRepo := repository.NewUserRepository(db)
	githubRepository := repository.NewGitHubReposRepository(db)

	// Services
	githubService := github.NewGithubService()
	userService := users.NewUserService(userRepo)
	reposService := repos.NewReposService(githubRepository)

	// Handlers
	userHandler := controller.NewUserHandler(userService)
	authHandler := controller.NewAuthHandler(userService, githubService)
	reposHandler := controller.NewReposHandler(reposService, githubService)

	// Versions
	v1 := app.Group("/api/v1")

	// Groups
	users := v1.Group("/users")
	auth := v1.Group("/auth")
	repos := v1.Group("/repos")

	// User routes
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUserByID)
	users.Get("/", userHandler.GetAllUsers)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Repos routes
	repos.Post("/sync", reposHandler.SyncRepositories)

	// Auth routes
	auth.Post("/sync", authHandler.Sync)

	return app
}
