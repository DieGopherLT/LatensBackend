package database

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrMongoURIMissing    = errors.New("mongo URI is required")
	ErrDatabaseClient     = errors.New("failed to create database client")
	ErrDatabaseConnection = errors.New("failed to connect to the database")
)

func Connect(uri, database string) (*mongo.Client, error) {

	if uri == "" {
		return nil, ErrMongoURIMissing
	}

	opts := options.Client().ApplyURI(uri).SetTimeout(5 * time.Second)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, ErrDatabaseClient
	}

	ctx, cancel := context.WithTimeout(context.Background(), *opts.Timeout)
	defer cancel()

	var result bson.M
	err = client.Database(database).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result)
	if err != nil {
		return nil, ErrDatabaseConnection
	}

	log.Default().Println("Connected to MongoDB")

	return client, nil
}
