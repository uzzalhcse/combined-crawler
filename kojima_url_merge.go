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

func processBatch(wg *sync.WaitGroup, client *mongo.Client, urls []string, batchSize int, productsCollection, productDetailsCollection *mongo.Collection) {
	defer wg.Done()

	// Prepare bulk operations for products and product_details collections
	var productOps, productDetailsOps []mongo.WriteModel

	for _, url := range urls {
		filter := bson.M{"url": url}
		update := bson.M{"$set": bson.M{"status": true}}

		productModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		productOps = append(productOps, productModel)

		productDetailsModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		productDetailsOps = append(productDetailsOps, productDetailsModel)

		// Execute the batch if it reaches the batchSize
		if len(productOps) >= batchSize {
			_, err := productsCollection.BulkWrite(context.TODO(), productOps)
			if err != nil {
				log.Printf("BulkWrite error (products): %v", err)
			}
			productOps = nil // Reset operations after executing the batch

			_, err = productDetailsCollection.BulkWrite(context.TODO(), productDetailsOps)
			if err != nil {
				log.Printf("BulkWrite error (product_details): %v", err)
			}
			productDetailsOps = nil // Reset operations after executing the batch
		}
	}

	// Execute remaining operations
	if len(productOps) > 0 {
		_, err := productsCollection.BulkWrite(context.TODO(), productOps)
		if err != nil {
			log.Printf("BulkWrite error (products): %v", err)
		}
	}

	if len(productDetailsOps) > 0 {
		_, err := productDetailsCollection.BulkWrite(context.TODO(), productDetailsOps)
		if err != nil {
			log.Printf("BulkWrite error (product_details): %v", err)
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

	// Select your database and collections
	database := client.Database("kojima_v2") // Replace with your database name
	productsCollection := database.Collection("products")
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
				go processBatch(&wg, client, urls, batchSize, productsCollection, productDetailsCollection)
				urls = nil // Reset urls slice after dispatching the batch
			}
		}
	}

	// Process any remaining URLs
	if len(urls) > 0 {
		wg.Add(1)
		go processBatch(&wg, client, urls, batchSize, productsCollection, productDetailsCollection)
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Product status updates completed.")
}
