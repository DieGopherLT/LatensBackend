package controller

import (
	"github.com/DieGopherLT/mfc_backend/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Handlers HTTP Request Handlers
type Handlers struct {
	Config   *config.Config
	Database *mongo.Client
}

// New initializes a new Handlers instance
func New(cfg *config.Config, db *mongo.Client) *Handlers {
	return &Handlers{
		Config:   cfg,
		Database: db,
	}
}
