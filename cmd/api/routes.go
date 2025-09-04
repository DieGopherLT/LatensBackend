package main

import (
	"github.com/DieGopherLT/mfc_backend/internal/controller"
	"github.com/DieGopherLT/mfc_backend/internal/database/repository"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupRoutes(app *fiber.App, db *mongo.Database) *fiber.App {

	userRepo := repository.NewUserRepository(db)
	userHandler := controller.NewUserHandler(userRepo)

	v1 := app.Group("/api/v1")
	users := v1.Group("/users")

	users.Post("/", userHandler.CreateUser)

	return app
}
