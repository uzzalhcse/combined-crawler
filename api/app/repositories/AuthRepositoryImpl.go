package repositories

import (
	"combined-crawler/api/app/models"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepositoryImpl struct {
	DB *mongo.Client
}

func NewAuthRepository(db *mongo.Client) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{DB: db}
}

func (r *AuthRepositoryImpl) FindUserByUsername(username string) (*models.User, error) {
	collection := r.DB.Database("your_database_name").Collection("users")

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) CreateUser(user *models.User) error {
	if r.DB == nil {
		return errors.New("DB is nil")
	}

	collection := r.DB.Database("your_database_name").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}

	return nil
}

func (r *AuthRepositoryImpl) UpdateUser(username string, updatedUser *models.User) error {
	if r.DB == nil {
		return errors.New("DB is nil")
	}

	collection := r.DB.Database("your_database_name").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	update := bson.M{"$set": updatedUser}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	return nil
}

func (r *AuthRepositoryImpl) FindUserByID(userID string) (*models.User, error) {
	collection := r.DB.Database("your_database_name").Collection("users")

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
