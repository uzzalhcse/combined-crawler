package generic_crawler

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

func loadSites(filename string) ([]SiteConfig, error) {
	var sites []SiteConfig

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&sites)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return sites, nil
}
func buildPlugin(path, functionName string) (string, error) {
	filename := "Plugin.go"
	src := filepath.Join(path, filename)
	pluginPath := filepath.Join(path, "bin", functionName+".so")

	if err := generatePlugin(src, pluginPath); err != nil {
		return "", fmt.Errorf("failed to build plugin %s: %v", functionName, err)
	}
	return pluginPath, nil
}
func generatePlugin(src, dest string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", dest, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
func getBaseUrl(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return ""
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host
	return baseURL
}
