package main

import (
	"fmt"
	"net/http"
	"strings"
)

func detectAutomation(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("User-Agent")
	fmt.Printf("Incoming request User-Agent: %s\n", userAgent)

	// Common headless browser indicators
	headlessIndicators := []string{
		"HeadlessChrome", // Chrome in headless mode
		"PhantomJS",      // PhantomJS (a headless WebKit scriptable with a JavaScript API)
		"Selenium",       // Selenium WebDriver
		"node-fetch",     // Node.js fetch API often used in automation
		"curl",           // curl request from terminal
		"Wget",           // Wget request from terminal
	}

	// Simple detection of common automation tools
	isAutomation := false
	for _, indicator := range headlessIndicators {
		if strings.Contains(userAgent, indicator) {
			isAutomation = true
			break
		}
	}

	// Additional detection based on lack of expected headers
	if r.Header.Get("Accept-Language") == "" || r.Header.Get("DNT") == "" {
		isAutomation = true
	}

	if isAutomation {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Automation tools detected. Access denied.")
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, real browser!")
	}
}

func main() {
	http.HandleFunc("/", detectAutomation)
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
