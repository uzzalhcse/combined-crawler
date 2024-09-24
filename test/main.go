package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
)

// Define a struct that matches the schema of your table
type PanasonicData struct {
	URL       string    `bigquery:"url"`
	HTMLData  string    `bigquery:"html_data"`
	CreatedAt time.Time `bigquery:"created_at"`
}

func main() {
	dataset := "panasonic_dataset"
	projectID := "lazuli-venturas-stg"
	table := "project_panasonic"
	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "test/gcp-file-upload-key.json"); err != nil {
		fmt.Printf("Failed to set GOOGLE_APPLICATION_CREDENTIALS: %v\n", err)
	}

	// Create context
	ctx := context.Background()

	// Initialize BigQuery client
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}
	defer client.Close()

	// Specify the dataset and table where data will be inserted
	inserter := client.Dataset(dataset).Table(table).Inserter()

	// Prepare the data to insert
	rows := []*PanasonicData{
		{
			URL:       "https://panasonic.jp/somepage.html",
			HTMLData:  "<html><body>Sample HTML data</body></html>",
			CreatedAt: time.Now(), // This field is used for partitioning
		},
	}

	// Insert data into the table
	if err := inserter.Put(ctx, rows); err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}

	log.Println("Data inserted successfully.")
}
