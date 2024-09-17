package ninjacrawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiEndpoint = "https://bq-relay-v2-beta-7tcydway2q-an.a.run.app"
	contentType = "application/json"
)

func (app *Crawler) submitProductData(productData *ProductDetail) error {
	jsonPayload, err := json.Marshal(productData)
	if err != nil {
		return fmt.Errorf("json conversion error: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiEndpoint+"/item/", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("%s: failed to create request: %w", productData.Url, err)
	}

	req.SetBasicAuth(app.Config.EnvString("API_USERNAME"), app.Config.EnvString("API_PASSWORD"))
	req.Header.Set("Content-Type", contentType)

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: failed to submit request: %w", productData.Url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: unexpected status code: %d and body: %s", productData.Url, response.StatusCode, response.Body)
	}

	return nil
}
