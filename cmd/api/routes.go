package main

import (
	"github.com/DieGopherLT/LatensBackend/internal/controller"
	"github.com/DieGopherLT/LatensBackend/internal/database/repository"
	"github.com/DieGopherLT/LatensBackend/internal/middleware"
	"github.com/DieGopherLT/LatensBackend/internal/services/github"
	"github.com/DieGopherLT/LatensBackend/internal/services/repos"
	"github.com/DieGopherLT/LatensBackend/internal/services/users"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupRoutes(app *fiber.App, db *mongo.Database) error {

	// Repositories
	userRepo := repository.NewUserRepository(db)
	githubRepository, err := repository.NewGitHubReposRepository(db)
	if err != nil {
		return err
	}

	// Services
	githubService := github.NewGithubService()
	userService := users.NewUserService(userRepo)
	reposService := repos.NewReposService(githubRepository)

	// Handlers
	userHandler := controller.NewUserHandler(userService)
	authHandler := controller.NewAuthHandler(userService, githubService)
	reposHandler := controller.NewReposHandler(reposService, githubService, userService)

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
	repos.Get("/", middleware.Guard(), reposHandler.GetRepos)
	repos.Post("/sync", middleware.Guard(), reposHandler.SyncRepos)

	// Auth routes
	auth.Post("/login", authHandler.Login)

	return nil
}
