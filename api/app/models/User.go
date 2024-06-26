package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `gorm:"unique" json:"email"`
	Password  string
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// UserIndexes returns the indexes for the User collection
func UserIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}}, // Create an index on the 'username' field
			Options: options.Index().SetUnique(true),
		},
	}
}
