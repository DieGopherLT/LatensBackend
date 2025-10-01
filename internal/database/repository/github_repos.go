package repository

import (
	"context"
	"errors"
	"time"

	"github.com/DieGopherLT/LatensBackend/internal/database/models"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoGitHubReposRepository struct {
	collection *mongo.Collection
}

func NewGitHubReposRepository(client *mongo.Database) (*MongoGitHubReposRepository, error) {
	collection := client.Collection("github_repositories")
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "github_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return nil, err
	}
	return &MongoGitHubReposRepository{collection: collection}, nil
}

func (r *MongoGitHubReposRepository) Create(ctx context.Context, repo *models.RepositoryDocument) error {
	filter := bson.M{"user_id": repo.UserID, "github_id": repo.GitHubID}
	update := bson.M{"$set": repo}
	opts := options.UpdateOne().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoGitHubReposRepository) CreateMany(ctx context.Context, repos []*models.RepositoryDocument) error {
	operations := lo.Map(repos, func(repo *models.RepositoryDocument, _ int) mongo.WriteModel {
		filter := bson.M{"user_id": repo.UserID, "github_id": repo.GitHubID}
		update := bson.M{"$set": repo}
		return mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
	})
	bulkOpts := options.BulkWrite().SetOrdered(false)
	_, err := r.collection.BulkWrite(ctx, operations, bulkOpts)
	return err
}

func (r *MongoGitHubReposRepository) FindByID(ctx context.Context, id string, userID string) (*models.RepositoryDocument, error) {
	var repo models.RepositoryDocument

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	condition := bson.M{"_id": objectID, "user_id": userObjectID}
	err = r.collection.FindOne(ctx, condition).Decode(&repo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *MongoGitHubReposRepository) FindAllByUser(ctx context.Context, userID string) ([]*models.RepositoryDocument, error) {
	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	condition := bson.M{"user_id": objectID}
	cursor, err := r.collection.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repos []*models.RepositoryDocument
	for cursor.Next(ctx) {
		var repo models.RepositoryDocument
		if err := cursor.Decode(&repo); err != nil {
			return nil, err
		}
		repos = append(repos, &repo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *MongoGitHubReposRepository) Update(ctx context.Context, id string, userID string, update map[string]any) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	condition := bson.M{"_id": objectID, "user_id": userObjectID}
	_, err = r.collection.UpdateOne(ctx, condition, bson.M{"$set": update})
	return err
}

func (r *MongoGitHubReposRepository) Delete(ctx context.Context, id string, userID string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	condition := bson.M{"_id": objectID, "user_id": userObjectID}
	_, err = r.collection.DeleteOne(ctx, condition)
	return err
}
