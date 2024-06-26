package bootstrap

import (
	"combined-crawler/api/app/exceptions"
	"combined-crawler/api/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createConnectionPool(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewDatabase(dbConfig config.DatabaseConfig) *mongo.Client {
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%s/",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
	)

	client, err := createConnectionPool(connString)
	if err != nil {
		exceptions.PanicIfNeeded(fmt.Errorf("[INIT] failed to connect to the database: %v", err))
	}

	fmt.Println("[INIT] Database connection established")

	//err = database.Migrate(client)
	//if err != nil {
	//	fmt.Println("DB Migration Error: ", err.Error())
	//}

	return client
}

func CloseDBConnection(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		exceptions.PanicIfNeeded(err)
	}
}
