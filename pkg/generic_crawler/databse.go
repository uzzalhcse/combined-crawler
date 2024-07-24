package generic_crawler

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

func (gc *GenericCrawler) mustGetClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	databaseURL := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		"lazuli",
		"x1RWo6cqFtHiaAHce5HB",
		"localhost",
		"27017",
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURL))
	if err != nil {
		panic(err)
	}

	// Check if the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB: %v", err)
		panic(err)
	}

	return client
}

// dropDatabase drops the specified database.
func (gc *GenericCrawler) dropDatabase() error {
	client := gc.mustGetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Database(gc.Name).Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

// // getCollection returns a collection from the database and ensures unique indexing.
func (gc *GenericCrawler) getCollection(collectionName string) *mongo.Collection {
	collection := gc.Database(gc.Name).Collection(collectionName)
	gc.ensureUniqueIndex(collection)
	return collection
}

// ensureUniqueIndex ensures that the "url" field in the collection has a unique index.
func (gc *GenericCrawler) ensureUniqueIndex(collection *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"url": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		fmt.Println("Could not create index: %v", err)
	}
}

// insert inserts multiple Url collections into the database.
func (gc *GenericCrawler) insert(model string, urlCollections []UrlCollection, parent string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var documents []interface{}
	for _, urlCollection := range urlCollections {
		urlCollection := UrlCollection{
			Url:       urlCollection.Url,
			Parent:    parent,
			Status:    false,
			Error:     false,
			MetaData:  urlCollection.MetaData,
			Attempts:  0,
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		}
		documents = append(documents, urlCollection)
	}

	opts := options.InsertMany().SetOrdered(false)
	collection := gc.getCollection(model)
	_, _ = collection.InsertMany(ctx, documents, opts)
}

// newSite creates a new site collection document in the database.
func (gc *GenericCrawler) newSite() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	document := SiteCollection{
		Url:       gc.Url,
		BaseUrl:   gc.BaseUrl,
		Status:    false,
		Attempts:  0,
		StartedAt: time.Now(),
		EndedAt:   nil,
	}

	collection := gc.getCollection(baseCollection)
	_, _ = collection.InsertOne(ctx, document)
}

// saveProductDetail saves or updates a product detail document in the database.

// markAsError marks a Url collection as having encountered an error and updates the database.
func (gc *GenericCrawler) markAsError(url string, dbCollection string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result bson.M

	collection := gc.getCollection(dbCollection)
	filter := bson.D{{Key: "url", Value: url}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return fmt.Errorf("Collection Not Found: %v", err)
	}
	timeNow := time.Now()
	attempts := result["attempts"].(int32) + 1
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "error", Value: true},
			{Key: "attempts", Value: attempts},
			{Key: "updated_at", Value: &timeNow},
		}},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[%s: => %s] could not mark as Error: Please check this [Error]: %v", dbCollection, url, err)
	}

	return nil
}

// markAsComplete marks a Url collection as having encountered an error and updates the database.
func (gc *GenericCrawler) markAsComplete(url string, dbCollection string) error {

	timeNow := time.Now()

	collection := gc.getCollection(dbCollection)

	filter := bson.D{{Key: "url", Value: url}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: true},
			{Key: "updated_at", Value: &timeNow},
		}},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("[:%s:%s] could not mark as Complete: Please check this [Error]: %v", dbCollection, url, err)
	}
	return nil
}
func (gc *GenericCrawler) SyncCurrentPageUrl(url, currentPageUrl string, dbCollection string) error {

	timeNow := time.Now()

	collection := gc.getCollection(dbCollection)

	filter := bson.D{{Key: "url", Value: url}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "current_page_url", Value: currentPageUrl},
			{Key: "updated_at", Value: &timeNow},
		}},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("[:%s:%s] could not sync currentpage [Error]: %v", dbCollection, url, err)
	}
	return nil
}

// getUrlsFromCollection retrieves URLs from a collection that meet specific criteria.
func (gc *GenericCrawler) getUrlsFromCollection(collection string) []string {
	filterCondition := bson.D{
		{Key: "status", Value: false},
		{Key: "attempts", Value: bson.D{{Key: "$lt", Value: 3}}},
	}
	return extractUrls(filterData(filterCondition, gc.getCollection(collection)))
}

// getUrlCollections retrieves Url collections from a collection that meet specific criteria.
func (gc *GenericCrawler) getUrlCollections(collection string) []UrlCollection {
	filterCondition := bson.D{
		{Key: "status", Value: false},
		{Key: "attempts", Value: bson.D{{Key: "$lt", Value: 3}}},
	}
	return gc.filterUrlData(filterCondition, gc.getCollection(collection))
}

// filterUrlData retrieves Url collections from a collection based on a filter condition.
func (gc *GenericCrawler) filterUrlData(filterCondition bson.D, mongoCollection *mongo.Collection) []UrlCollection {
	findOptions := options.Find().SetLimit(1000)

	cursor, err := mongoCollection.Find(context.TODO(), filterCondition, findOptions)
	if err != nil {
		fmt.Println(err.Error())
	}

	var results []UrlCollection
	if err = cursor.All(context.TODO(), &results); err != nil {
		fmt.Println(err.Error())
	}

	return results
}

// filterData retrieves documents from a collection based on a filter condition.
func filterData(filterCondition bson.D, mongoCollection *mongo.Collection) []bson.M {
	findOptions := options.Find().SetLimit(1000)

	cursor, err := mongoCollection.Find(context.TODO(), filterCondition, findOptions)
	if err != nil {
		slog.Error(err.Error())
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		slog.Error(err.Error())
	}

	return results
}

// extractUrls extracts URLs from a list of BSON documents.
func extractUrls(results []bson.M) []string {
	var urls []string
	for _, result := range results {
		if url, ok := result["url"].(string); ok {
			urls = append(urls, url)
		}
	}
	return urls
}

// close closes the MongoDB client connection.
func (gc *GenericCrawler) closeClient() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return gc.Disconnect(ctx)
}
