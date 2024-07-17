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
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sourcePath := filepath.Join(cwd)
	outputPath := filepath.Join(cwd, "dist", "crawler")
	cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("Error building site: %v\nOutput: %s", err, output)
	}
	GCP_CREDENTIALS_PATH := "gcp-file-upload-key.json"
	fmt.Println("GCP_CREDENTIALS_PATH: ", GCP_CREDENTIALS_PATH)
	err = ninjacrawler.UploadToGCPBucket("crawler", GCP_CREDENTIALS_PATH, outputPath, "crawler")
	if err != nil {
		log.Fatalf("Error uploading built site to bucket: %v", err)
		return
	}

	fmt.Println("Build successful ", string(output))
}
