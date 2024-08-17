package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// Define the batch size
	batchSize := 1000

	// Set up options for the cursor
	findOptions := options.Find().SetBatchSize(int32(batchSize))

	// Find all documents in the source collection with the specified options
	cursor, err := sourceCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var documents []interface{}
	var totalCopied int

	// Iterate over the cursor and decode each document in batches
	for cursor.Next(context.TODO()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			log.Fatal(err)
		}
		documents = append(documents, document)

		// When the batch is full, insert it into the destination collection
		if len(documents) >= batchSize {
			insertResult, err := destinationCollection.InsertMany(context.TODO(), documents)
			if err != nil {
				log.Fatal(err)
			}
			totalCopied += len(insertResult.InsertedIDs)
			fmt.Printf("Copied %d documents so far...\n", totalCopied)
			documents = nil // Clear the slice for the next batch

			// Sleep to prevent overwhelming the server
			time.Sleep(1 * time.Second)
		}
	}

	// Insert any remaining documents that didn't fill a full batch
	if len(documents) > 0 {
		insertResult, err := destinationCollection.InsertMany(context.TODO(), documents)
		if err != nil {
			log.Fatal(err)
		}
		totalCopied += len(insertResult.InsertedIDs)
	}

	fmt.Printf("Copied %d documents to the destination collection\n", totalCopied)
}
