package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
)

func processBatch(wg *sync.WaitGroup, client *mongo.Client, urls []string, batchSize int, productsCollection, productDetailsCollection *mongo.Collection, processedCount *int32, total int32) {
	defer wg.Done()

	var productOps, productDetailsOps []mongo.WriteModel

	for _, url := range urls {
		filter := bson.M{"url": url}
		update := bson.M{"$set": bson.M{"status": true}}

		productModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		productOps = append(productOps, productModel)

		productDetailsModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		productDetailsOps = append(productDetailsOps, productDetailsModel)

		if len(productOps) >= batchSize {
			_, err := productsCollection.BulkWrite(context.TODO(), productOps)
			if err != nil {
				log.Printf("BulkWrite error (products): %v", err)
			}
			productOps = nil

			_, err = productDetailsCollection.BulkWrite(context.TODO(), productDetailsOps)
			if err != nil {
				log.Printf("BulkWrite error (product_details): %v", err)
			}
			productDetailsOps = nil
		}
	}

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

	atomic.AddInt32(processedCount, int32(len(urls)))
	progress := float32(*processedCount) / float32(total) * 100
	fmt.Printf("Progress: %.2f%%\n", progress)
}

func main() {
	uri := "mongodb://lazuli:x1RWo6cqFtHiaAHce5HB@localhost:27017/"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	database := client.Database("kojima_v2")
	productsCollection := database.Collection("products")
	productDetailsCollection := database.Collection("product_details")

	total, err := productDetailsCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total documents to process: %d\n", total)

	findOptions := options.Find().SetBatchSize(500) // Set a smaller batch size to avoid CursorNotFound error
	cursor, err := productDetailsCollection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var urls []string
	batchSize := 1000
	numWorkers := runtime.NumCPU()
	var processedCount int32

	var wg sync.WaitGroup

	for cursor.Next(context.TODO()) {
		var productDetail bson.M
		if err := cursor.Decode(&productDetail); err != nil {
			log.Fatal(err)
		}

		if productDetailURL, ok := productDetail["url"].(string); ok {
			urls = append(urls, productDetailURL)

			if len(urls) >= batchSize*numWorkers {
				wg.Add(1)
				go processBatch(&wg, client, urls, batchSize, productsCollection, productDetailsCollection, &processedCount, int32(total))
				urls = nil
			}
		}
	}

	if len(urls) > 0 {
		wg.Add(1)
		go processBatch(&wg, client, urls, batchSize, productsCollection, productDetailsCollection, &processedCount, int32(total))
	}

	wg.Wait()

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Product status updates completed.")
}
