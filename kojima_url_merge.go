package main

import (
	"context"
	"log"
	"runtime"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func processBatch(wg *sync.WaitGroup, client *mongo.Client, urls []string, batchSize int) {
	defer wg.Done()

	// Select the database and collection
	database := client.Database("kojima_test") // Replace with your database name
	productsCollection := database.Collection("products")

	// Prepare bulk operations
	var operations []mongo.WriteModel

	for _, url := range urls {
		filter := bson.M{"url": url}
		update := bson.M{"$set": bson.M{"status": true}}

		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		operations = append(operations, model)

		// Execute the batch if it reaches the batchSize
		if len(operations) >= batchSize {
			_, err := productsCollection.BulkWrite(context.TODO(), operations)
			if err != nil {
				log.Printf("BulkWrite error: %v", err)
			}
			operations = nil // Reset operations after executing the batch
		}
	}

	// Execute remaining operations
	if len(operations) > 0 {
		_, err := productsCollection.BulkWrite(context.TODO(), operations)
		if err != nil {
			log.Printf("BulkWrite error: %v", err)
		}
	}
}

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
	database := client.Database("kojima_test") // Replace with your database name
	productDetailsCollection := database.Collection("product_details")

	// Find all product_details documents and retrieve URLs
	cursor, err := productDetailsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var urls []string
	batchSize := 1000              // Define the batch size for updates
	numWorkers := runtime.NumCPU() // Number of parallel workers

	// WaitGroup to synchronize Goroutines
	var wg sync.WaitGroup

	for cursor.Next(context.TODO()) {
		var productDetail bson.M
		if err := cursor.Decode(&productDetail); err != nil {
			log.Fatal(err)
		}

		// Extract the URL from product_details
		if productDetailURL, ok := productDetail["url"].(string); ok {
			urls = append(urls, productDetailURL)

			// When the batch is full, process it
			if len(urls) >= batchSize*numWorkers {
				wg.Add(1)
				go processBatch(&wg, client, urls, batchSize)
				urls = nil // Reset urls slice after dispatching the batch
			}
		}
	}

	// Process any remaining URLs
	if len(urls) > 0 {
		wg.Add(1)
		go processBatch(&wg, client, urls, batchSize)
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Product status updates completed.")
}
