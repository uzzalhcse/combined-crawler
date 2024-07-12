package ninjacrawler

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// CrawlPageDetail initiates the crawling process for detailed page information from the specified collection.
// It distributes the work among multiple goroutines and uses proxies if available.
func (app *Crawler) CrawlPageDetail(processorConfigs []ProcessorConfig) {
	for _, processorConfig := range processorConfigs {
		overrideEngineDefaults(app.engine, &processorConfig.Engine)
		app.toggleClient()
		processedUrls := make(map[string]bool) // Track processed URLs
		total := int32(0)
		app.crawlPageDetailRecursive(processorConfig, processedUrls, &total, 0)
		if atomic.LoadInt32(&total) > 0 {
			app.Logger.Info("Total %v %v Inserted ", atomic.LoadInt32(&total), processorConfig.OriginCollection)
		}
		exportProductDetailsToCSV(app, processorConfig.Entity, 1)
	}
}

func (app *Crawler) crawlPageDetailRecursive(processorConfig ProcessorConfig, processedUrls map[string]bool, total *int32, counter int32) {
	urlCollections := app.getUrlCollections(processorConfig.OriginCollection)
	newUrlCollections := []UrlCollection{}
	for i, urlCollection := range urlCollections {
		if app.isLocalEnv && i >= app.engine.DevCrawlLimit {
			break
		}
		if !processedUrls[urlCollection.CurrentPageUrl] || !processedUrls[urlCollection.Url] {
			newUrlCollections = append(newUrlCollections, urlCollection)
		}
	}
	var wg sync.WaitGroup
	urlChan := make(chan UrlCollection, len(newUrlCollections))
	resultChan := make(chan interface{}, len(newUrlCollections))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, urlCollection := range newUrlCollections {
		urlChan <- urlCollection
		if urlCollection.Attempts > 0 && urlCollection.Attempts <= app.engine.MaxRetryAttempts {
			processedUrls[urlCollection.CurrentPageUrl] = false // Do Not Mark URL as processed
			processedUrls[urlCollection.Url] = false            // Do Not Mark URL as processed
		} else {
			processedUrls[urlCollection.CurrentPageUrl] = true // Mark URL as processed
			processedUrls[urlCollection.Url] = true            // Mark URL as processed
		}
	}
	close(urlChan)

	proxyCount := len(app.engine.ProxyServers)
	batchSize := app.engine.ConcurrentLimit
	totalUrls := len(newUrlCollections)
	goroutineCount := min(max(proxyCount, 1)*batchSize, totalUrls) // Determine the required number of goroutines

	for i := 0; i < goroutineCount; i++ {
		proxy := Proxy{}
		if proxyCount > 0 {
			proxy = app.engine.ProxyServers[i%proxyCount]
		}
		wg.Add(1)
		go func(proxy Proxy) {
			defer wg.Done()
			app.crawlWorker(ctx, processorConfig, urlChan, resultChan, proxy, app.isLocalEnv, &counter)
		}(proxy)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for results := range resultChan {
		switch v := results.(type) {
		case CrawlResult:
			switch res := v.Results.(type) {
			case *ProductDetail:
				err := app.handleProductDetail(res, processorConfig, v)
				if err != nil {
					app.Logger.Error(err.Error())
					continue
				}
				atomic.AddInt32(total, 1)
			}
		}
	}
	if app.isLocalEnv && atomic.LoadInt32(&counter) >= int32(app.engine.DevCrawlLimit) {
		cancel()
		return
	}
	if len(newUrlCollections) > 0 {
		app.crawlPageDetailRecursive(processorConfig, processedUrls, total, counter)
	}
}

// validateRequiredFields checks if the required fields are non-empty in the ProductDetail struct.
// Returns two slices: one for invalid fields and one for unknown fields.
func validateRequiredFields(product *ProductDetail, validationRules []string) ([]string, []string) {
	var invalidFields []string
	var unknownFields []string

	v := reflect.ValueOf(*product)
	t := v.Type()

	for _, rule := range validationRules {
		parts := strings.Split(rule, "|")
		field := parts[0]
		rules := parts[1:]

		f, ok := t.FieldByName(field)
		if !ok {
			unknownFields = append(unknownFields, field)
			continue
		}

		fieldValue := v.FieldByName(field)
		fieldValueStr := fmt.Sprintf("%v", fieldValue.Interface())
		if fieldValueStr == "" {
			invalidFields = append(invalidFields, fmt.Sprintf("%s: required", f.Name))
		}
		for _, r := range rules {
			ruleParts := strings.SplitN(r, ":", 2)
			ruleName := ruleParts[0]
			ruleValue := ""
			if len(ruleParts) > 1 {
				ruleValue = ruleParts[1]
			}

			switch ruleName {
			case "required":
				if fieldValueStr == "" {
					invalidFields = append(invalidFields, fmt.Sprintf("%s: required", f.Name))
				}
			case "string":
				if fieldValue.Kind() != reflect.String {
					invalidFields = append(invalidFields, fmt.Sprintf("%s: not a string", f.Name))
				}
			case "max":
				maxLength, err := strconv.Atoi(ruleValue)
				if err == nil && len(fieldValueStr) > maxLength {
					invalidFields = append(invalidFields, fmt.Sprintf("%s: exceeds max length of %d", f.Name, maxLength))
				}
			case "trim":
				if strings.TrimSpace(fieldValueStr) != fieldValueStr {
					invalidFields = append(invalidFields, fmt.Sprintf("%s: not trimmed", f.Name))
				}
			case "blacklists":
				excludeValues := strings.Split(ruleValue, ",")
				for _, excludeValue := range excludeValues {
					excludeValue = strings.TrimSpace(excludeValue)
					if fieldValueStr == excludeValue {
						invalidFields = append(invalidFields, fmt.Sprintf("%s: blacklist value '%s'", f.Name, excludeValue))
						break
					}
				}
			// Add more cases for other validation rules as needed
			default:
			}
		}
	}
	return invalidFields, unknownFields
}

func (app *Crawler) handleProductDetail(res *ProductDetail, processorConfig ProcessorConfig, v CrawlResult) error {
	invalidFields, unknownFields := validateRequiredFields(res, processorConfig.Preference.ValidationRules)
	if len(unknownFields) > 0 {
		return fmt.Errorf("unknown fields provided: %v", unknownFields)
	}
	if len(invalidFields) > 0 {
		msg := fmt.Sprintf("Validation failed: %v\n", invalidFields)
		html, _ := v.Document.Html()
		if app.engine.IsDynamic {
			html = app.getHtmlFromPage(v.Page)
		}
		app.Logger.Html(html, v.UrlCollection.Url, msg)
		err := app.markAsError(v.UrlCollection.Url, processorConfig.OriginCollection)
		if err != nil {
			return err
		}
	}

	app.saveProductDetail(processorConfig.Entity, res)
	if !app.isLocalEnv {
		err := app.submitProductData(res)
		if err != nil {
			app.Logger.Fatal("Failed to submit product data to API Server: %v", err)
			err := app.markAsError(v.UrlCollection.Url, processorConfig.OriginCollection)
			if err != nil {
				return err
			}
		}
	}

	if !processorConfig.Preference.DoNotMarkAsComplete {
		err := app.markAsComplete(v.UrlCollection.Url, processorConfig.OriginCollection)
		if err != nil {
			return err
		}
	}
	return nil
}
