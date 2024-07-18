package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GenerateBinaryBuild(SiteID string) error {

	//appsDir := "/root/ninja-combined-crawler/apps"
	//distDir := "/root/combined-crawler"

	appsDir := "/home/uzzal/Workplace/Lazuli/ninja-combined-crawler/apps"
	distDir := "/home/uzzal/Workplace/github/combined-crawler"

	// Get the absolute path of the parent directory
	parentDir, err := filepath.Abs(distDir)
	if err != nil {
		return fmt.Errorf("Error getting parent directory: %v", err)
	}

	files, err := os.ReadDir(appsDir)
	if err != nil {
		return fmt.Errorf("Error reading directory:", err)
	}

	siteFound := false
	for _, file := range files {
		if file.IsDir() {
			dirname := file.Name()
			if SiteID == dirname {
				siteFound = true
				fmt.Printf("Generating Binary for: %s\n", dirname)
				outputPath := filepath.Join(parentDir, "dist", dirname)
				sourcePath := fmt.Sprintf("%s/%s", appsDir, dirname)
				fmt.Println("sourcePath: ", sourcePath)
				fmt.Println("outputPath: ", outputPath)
				cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && git pull && go build -o %s", sourcePath, outputPath))
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("Error building site: %v\nOutput: %s", err, output)
				}
			}
		}
	}
	if !siteFound {
		return fmt.Errorf("invalid site: %s", SiteID)
	}
	return nil
}

func CreateVM(siteID, zone string) (map[string]interface{}, error) {
	projectID := "lazuli-venturas"
	instanceName := siteID

	// Get gcloud access token
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error retrieving access token: %v", err)
	}
	accessToken := strings.TrimSpace(string(output))

	// Construct the request body for creating the VM
	vmRequestBody := map[string]interface{}{
		"canIpForward":       false,
		"deletionProtection": false,
		"description":        "",
		"disks": []map[string]interface{}{
			{
				"autoDelete": true,
				"boot":       true,
				"deviceName": instanceName,
				"initializeParams": map[string]interface{}{
					"diskSizeGb":  "10",
					"diskType":    "projects/lazuli-venturas/zones/asia-northeast1-c/diskTypes/pd-balanced",
					"sourceImage": "projects/ubuntu-os-cloud/global/images/ubuntu-2204-jammy-v20240701",
				},
				"mode": "READ_WRITE",
				"type": "PERSISTENT",
			},
		},
		"displayDevice": map[string]bool{
			"enableDisplay": false,
		},
		"guestAccelerators":       []interface{}{},
		"instanceEncryptionKey":   map[string]string{},
		"keyRevocationActionType": "NONE",
		"labels": map[string]string{
			"goog-ec-src": "vm_add-rest",
		},
		"machineType": "projects/lazuli-venturas/zones/asia-northeast1-a/machineTypes/e2-medium",
		"metadata": map[string]interface{}{
			"items": []interface{}{},
		},
		"name": instanceName,
		"networkInterfaces": []map[string]interface{}{
			{
				"accessConfigs": []map[string]string{
					{
						"name":        "External NAT",
						"networkTier": "PREMIUM",
					},
				},
				"stackType":  "IPV4_ONLY",
				"subnetwork": "projects/lazuli-venturas/regions/asia-northeast1/subnetworks/default",
			},
		},
		"params": map[string]interface{}{
			"resourceManagerTags": map[string]string{},
		},
		"reservationAffinity": map[string]string{
			"consumeReservationType": "ANY_RESERVATION",
		},
		"scheduling": map[string]interface{}{
			"automaticRestart":  true,
			"onHostMaintenance": "MIGRATE",
			"provisioningModel": "STANDARD",
		},
		"serviceAccounts": []map[string]interface{}{
			{
				"email": "845643578999-compute@developer.gserviceaccount.com",
				"scopes": []string{
					"https://www.googleapis.com/auth/devstorage.read_only",
					"https://www.googleapis.com/auth/logging.write",
					"https://www.googleapis.com/auth/monitoring.write",
					"https://www.googleapis.com/auth/service.management.readonly",
					"https://www.googleapis.com/auth/servicecontrol",
					"https://www.googleapis.com/auth/trace.append",
				},
			},
		},
		"shieldedInstanceConfig": map[string]interface{}{
			"enableIntegrityMonitoring": true,
			"enableSecureBoot":          false,
			"enableVtpm":                true,
		},
		"tags": map[string]interface{}{
			"items": []string{"http-server", "https-server"},
		},
		"zone": fmt.Sprintf("projects/%s/zones/%s", projectID, zone),
	}

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(vmRequestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON request body: %v", err)
	}

	// Send the request to create the VM
	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/instances", projectID, zone)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to create instance: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return jsonResponse, nil
}
