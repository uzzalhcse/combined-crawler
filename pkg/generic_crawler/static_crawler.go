package generic_crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"strings"
	"time"
)

func (gc *GenericCrawler) GetHttpClient() *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return client
}

func (gc *GenericCrawler) NavigateToStaticURL(urlString string) (*goquery.Document, error) {
	body, err := gc.getResponseBody(gc.httpClient, urlString)
	if err != nil {
		return nil, err
	}

	// Create a reader that can decode the response body with the correct encoding
	reader, err := charset.NewReader(strings.NewReader(string(body)), "")
	if err != nil {
		return nil, fmt.Errorf("failed to create reader with correct encoding: %w", err)
	}

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	return document, nil
}

func (gc *GenericCrawler) getResponseBody(client *http.Client, urlString string) ([]byte, error) {
	proxyIp := ""

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: from %s to %v", proxyIp, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("failed to fetch page: %v", resp.Status)
		return nil, fmt.Errorf(msg)
	}
	return body, nil
}
