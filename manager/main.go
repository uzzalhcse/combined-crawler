package main

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	sitesDir := "/root/ninja-combined-crawler/apps"

	// Get the absolute path of the parent directory
	parentDir, err := filepath.Abs("/root/combined-crawler")
	if err != nil {
		log.Fatalf("Error getting parent directory: %v", err)
	}

	files, err := os.ReadDir(sitesDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			dirname := file.Name()

			fmt.Printf("starting vm %s\n", dirname)
			outputPath := filepath.Join(parentDir, "dist", dirname)
			fmt.Println("outputPath: ", outputPath)
			cmd := exec.Command("go", "build", "-o", outputPath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Error building site: %v\nOutput: %s", err, output)
			}
		}
	}
}

func buildAndUpload() {

	// Get the absolute path of the parent directory
	parentDir, err := filepath.Abs("./../")
	if err != nil {
		log.Fatalf("Error getting parent directory: %v", err)
	}

	// Set the source and output paths
	sourcePath := parentDir
	outputPath := filepath.Join(parentDir, "dist", "crawler")

	// Build the project
	cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error building site: %v\nOutput: %s", err, output)
	}

	// Set the GCP credentials path
	gcpCredentialsPath := filepath.Join(parentDir, "gcp-file-upload-key.json")
	// Upload the built site to the GCP bucket
	err = ninjacrawler.UploadToGCPBucket("crawler", gcpCredentialsPath, outputPath, "crawler")
	if err != nil {
		log.Fatalf("Error uploading built site to bucket: %v", err)
		return
	}

	fmt.Println("Build successful:", string(output))
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
