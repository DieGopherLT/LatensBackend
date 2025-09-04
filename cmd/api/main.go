package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DieGopherLT/mfc_backend/internal/config"
	"github.com/DieGopherLT/mfc_backend/internal/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	mongo, db, err := database.Connect(cfg.MongoDbURI, cfg.MongoDbName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = mongo.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()

	setupRoutes(app, db)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	<-c
	log.Println("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatal("Failed to shutdown server: ", err)
	}
}
