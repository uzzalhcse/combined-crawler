package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	// Select the source and destination databases
	sourceDatabase := client.Database("kojima")         // Replace with your source database name
	destinationDatabase := client.Database("kojima_v2") // Replace with your destination database name

	// Select the collections
	sourceCollection := sourceDatabase.Collection("product_details")
	destinationCollection := destinationDatabase.Collection("product_details")

	// Find all documents in the source collection
	cursor, err := sourceCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var documents []interface{}

	// Iterate over the cursor and decode each document
	for cursor.Next(context.TODO()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			log.Fatal(err)
		}
		documents = append(documents, document)
	}

	// Insert all documents into the destination collection
	if len(documents) > 0 {
		insertResult, err := destinationCollection.InsertMany(context.TODO(), documents)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Copied %d documents to the destination collection\n", len(insertResult.InsertedIDs))
	} else {
		fmt.Println("No documents found in the source collection.")
	}
}
