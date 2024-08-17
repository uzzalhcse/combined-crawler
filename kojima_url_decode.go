package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// MongoDB connection URI
	uri := "mongodb://lazuli:x1RWo6cqFtHiaAHce5HB@localhost:27017/" // Replace with your MongoDB URI

	// Create a MongoDB client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Ping the database to check the connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	// Select your database and collection
	database := client.Database("kojima")                // Replace with your database name
	collection := database.Collection("product_details") // Replace with your collection name

	// Filter to find documents with URL-encoded fields (optional, can be removed if updating all documents)
	filter := bson.M{
		"url": bson.M{"$exists": true},
	}

	// Find documents
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			log.Fatal(err)
		}

		// Decode the URL fields
		decodedURL, _ := url.QueryUnescape(document["url"].(string))

		// Update the document with the decoded URLs
		update := bson.M{
			"$set": bson.M{
				"url": decodedURL,
			},
		}

		// Perform the update
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": document["_id"]}, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Updated document with _id: %v\n", document["_id"])
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("URL decoding and update completed.")
}
