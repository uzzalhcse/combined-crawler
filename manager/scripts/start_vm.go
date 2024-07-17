package manager

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func StartVm() {
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
		fmt.Printf("Instance %s not found.\n", instanceName)
		os.Exit(1)
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
