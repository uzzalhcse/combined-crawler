package database

import (
	"combined-crawler/api/app/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func Migrate(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Example: Ensure the 'testmodels' collection exists and has an index
	err := createCollectionAndIndexes(ctx, client, "ninja_crawler", "testmodels", models.TestModelIndexes())
	if err != nil {
		return err
	}

	// Example: Ensure the 'users' collection exists and has an index
	err = createCollectionAndIndexes(ctx, client, "ninja_crawler", "users", models.UserIndexes())
	if err != nil {
		return err
	}

	return nil
}

func createCollectionAndIndexes(ctx context.Context, client *mongo.Client, dbName, collectionName string, indexes []mongo.IndexModel) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Create indexes
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	return nil
}
