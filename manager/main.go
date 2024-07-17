package main

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

func main() {
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
