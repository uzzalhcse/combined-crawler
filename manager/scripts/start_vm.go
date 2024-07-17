package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Usage: go run start_vm.go INSTANCE_NAME")
		os.Exit(1)
	}

	instanceName := args[0]
	projectID := "lazuli-venturas"

	// Get gcloud access token
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error retrieving access token: %v", err)
	}
	accessToken := strings.TrimSpace(string(output))

	// Retrieve the zone of the instance
	cmd = exec.Command("gcloud", "compute", "instances", "list", "--filter=name="+instanceName, "--format=value(zone)")
	output, err = cmd.Output()
	if err != nil {
		log.Fatalf("Error retrieving instance zone: %v", err)
	}
	zone := strings.TrimSpace(string(output))
	fmt.Printf("Instance located at %s Zone.\n", zone)
	if zone == "" {
		fmt.Printf("Instance %s not found. Creating it...\n", instanceName)
		err := createVM(instanceName, projectID, accessToken)
		if err != nil {
			log.Fatalf("Error creating instance: %v", err)
		}
		fmt.Printf("Instance %s has been created.\n", instanceName)
		return
	}

	// Get the status of the VM
	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s", projectID, zone, instanceName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	var instanceData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&instanceData)
	if err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	status := instanceData["status"].(string)
	if status != "RUNNING" {
		// Start the VM if it is not running
		startURL := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s/start", projectID, zone, instanceName)
		startReq, err := http.NewRequest("POST", startURL, nil)
		if err != nil {
			log.Fatalf("Error creating start request: %v", err)
		}
		startReq.Header.Set("Authorization", "Bearer "+accessToken)

		startResp, err := client.Do(startReq)
		if err != nil {
			log.Fatalf("Error starting instance: %v", err)
		}
		defer startResp.Body.Close()

		fmt.Printf("Instance %s in zone %s has been started.\n", instanceName, zone)
	} else {
		fmt.Printf("Instance %s in zone %s is already running.\n", instanceName, zone)
	}
}

func createVM(instanceName, projectID, accessToken string) error {
	zone := "asia-northeast1-a"

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
		return fmt.Errorf("error marshaling JSON request body: %v", err)
	}

	// Send the request to create the VM
	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/instances", projectID, zone)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to create instance: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}
