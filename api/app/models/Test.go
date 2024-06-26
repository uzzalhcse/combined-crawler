package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type TestModel struct {
	gorm.Model
	Name  string
	Email string
	// Add your fields here
}

// TestModelIndexes returns the indexes for the TestModel collection
func TestModelIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}}, // Create an index on the 'name' field
			Options: options.Index().SetUnique(true),
		},
	}
}
