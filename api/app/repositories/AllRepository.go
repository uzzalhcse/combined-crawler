package repositories

import (
	"combined-crawler/api/app/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const DBName = "crawl_manager"

type Repository struct {
	DB *mongo.Client
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{DB: db}
}

// SiteCollection CRUD
var siteCollection models.SiteCollection

func (r *Repository) GetAllSiteCollections() ([]models.SiteCollection, error) {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.SiteCollection
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (r *Repository) CreateSiteCollection(siteCollection *models.SiteCollection) error {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, siteCollection)
	return err
}

func (r *Repository) GetSiteCollectionByID(siteID string) (*models.SiteCollection, error) {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var siteCollection models.SiteCollection
	err := collection.FindOne(ctx, bson.M{"site_id": siteID}).Decode(&siteCollection)
	if err != nil {
		return nil, err
	}
	return &siteCollection, nil
}

func (r *Repository) UpdateSiteCollection(siteID string, update bson.M) error {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"site_id": siteID}, bson.M{"$set": update})
	return err
}

func (r *Repository) DeleteSiteCollection(siteID string) error {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"site_id": siteID})
	return err
}

// Collection CRUD
var collection models.Collection

func (r *Repository) GetAllCollections() ([]models.Collection, error) {
	collection := r.DB.Database(DBName).Collection(collection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Collection
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (r *Repository) CreateCollection(collection *models.Collection) error {
	collectionColl := r.DB.Database(DBName).Collection(collection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collectionColl.InsertOne(ctx, collection)
	return err
}

func (r *Repository) GetCollectionByID(collectionID string) (*models.Collection, error) {
	collectionColl := r.DB.Database(DBName).Collection(collection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var collection models.Collection
	err := collectionColl.FindOne(ctx, bson.M{"collection_id": collectionID}).Decode(&collection)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *Repository) UpdateCollection(collectionID string, update bson.M) error {
	collectionColl := r.DB.Database(DBName).Collection(collection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collectionColl.UpdateOne(ctx, bson.M{"collection_id": collectionID}, bson.M{"$set": update})
	return err
}

func (r *Repository) DeleteCollection(collectionID string) error {
	collectionColl := r.DB.Database(DBName).Collection(collection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collectionColl.DeleteOne(ctx, bson.M{"collection_id": collectionID})
	return err
}

// UrlCollection CRUD
var urlCollection models.UrlCollection

func (r *Repository) CreateUrlCollection(urlCollection *models.UrlCollection) error {
	urlCollectionColl := r.DB.Database(DBName).Collection(urlCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := urlCollectionColl.InsertOne(ctx, urlCollection)
	return err
}

func (r *Repository) GetUrlCollectionByID(collectionID string) (*models.UrlCollection, error) {
	urlCollectionColl := r.DB.Database(DBName).Collection(urlCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var urlCollection models.UrlCollection
	err := urlCollectionColl.FindOne(ctx, bson.M{"collection_id": collectionID}).Decode(&urlCollection)
	if err != nil {
		return nil, err
	}
	return &urlCollection, nil
}

func (r *Repository) UpdateUrlCollection(collectionID string, update bson.M) error {
	urlCollectionColl := r.DB.Database(DBName).Collection(urlCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := urlCollectionColl.UpdateOne(ctx, bson.M{"collection_id": collectionID}, bson.M{"$set": update})
	return err
}

func (r *Repository) DeleteUrlCollection(collectionID string) error {
	urlCollectionColl := r.DB.Database(DBName).Collection(urlCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := urlCollectionColl.DeleteOne(ctx, bson.M{"collection_id": collectionID})
	return err
}

var siteSecretCollection models.SiteSecret

func (r *Repository) CreateSecretCollection(siteSecret *models.SiteSecret) error {
	collection := r.DB.Database(DBName).Collection(siteSecret.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, siteSecret)
	return err
}
func (r *Repository) GetAllSiteSecretCollections() ([]models.SiteSecret, error) {
	collection := r.DB.Database(DBName).Collection(siteSecretCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.SiteSecret
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (r *Repository) GetSiteSecretCollectionByID(siteID string) (*models.SiteSecret, error) {
	collection := r.DB.Database(DBName).Collection(siteSecretCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var secretCollection models.SiteSecret
	err := collection.FindOne(ctx, bson.M{"site_id": siteID}).Decode(&secretCollection)
	if err != nil {
		return nil, err
	}
	return &secretCollection, nil
}

var globalSecretCollection models.GlobalSecret

func (r *Repository) GetAllGlobalSecretCollections() ([]models.GlobalSecret, error) {
	collection := r.DB.Database(DBName).Collection(globalSecretCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.GlobalSecret
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
