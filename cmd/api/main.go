package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DieGopherLT/mfc_backend/internal/config"
	"github.com/DieGopherLT/mfc_backend/internal/controller"
	"github.com/DieGopherLT/mfc_backend/internal/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.MongoDbURI, cfg.MongoDbName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()
	handlers := controller.New(cfg, db)

	setupRoutes(app, handlers)

	app.Listen(fmt.Sprintf(":%s", cfg.Port))
}
