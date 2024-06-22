package ninjacrawler

import (
	"context"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"reflect"
	"sync"
	"sync/atomic"
)

// Struct to hold both results and the UrlCollection
type CrawlResult struct {
	Results       interface{}
	UrlCollection UrlCollection
	Page          playwright.Page
}

func (app *Crawler) crawlWorker(ctx context.Context, dbCollection string, urlChan <-chan UrlCollection, resultChan chan<- interface{}, proxy Proxy, processor interface{}, isLocalEnv bool, counter *int32) {
	browser, page, err := app.GetBrowserPage(app.pw, app.engine.BrowserType, proxy)
	if err != nil {
		app.Logger.Fatal("failed to initialize browser with Proxy: %v\n", err)
	}
	defer browser.Close()
	defer page.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case urlCollection, more := <-urlChan:
			if !more {
				return
			}

			if isLocalEnv && atomic.LoadInt32(counter) >= int32(app.engine.DevCrawlLimit) {
				app.Logger.Warn("Dev Crawl limit reached!")
				return
			}

			if proxy.Server != "" {
				app.Logger.Info("Crawling :%s: %s using Proxy %s", dbCollection, urlCollection.Url, proxy.Server)
			} else {
				app.Logger.Info("Crawling :%s: %s", dbCollection, urlCollection.Url)
			}

			doc, err := app.NavigateToURL(page, urlCollection.Url)
			if err != nil {
				err := app.markAsError(urlCollection.Url, dbCollection)
				if err != nil {
					app.Logger.Info(err.Error())
					return
				}
				continue
			}

			crawlerCtx := CrawlerContext{
				App:           app,
				Document:      doc,
				UrlCollection: urlCollection,
				Page:          page,
			}

			var results interface{}
			switch v := processor.(type) {
			case func(CrawlerContext) []UrlCollection:
				results = v(crawlerCtx)

			case UrlSelector:
				results = app.processDocument(doc, v, urlCollection)

			case ProductDetailSelector:
				results = crawlerCtx.handleProductDetail()

			default:
				app.Logger.Fatal("Unsupported processor type: %T", processor)
			}

			crawlResult := CrawlResult{
				Results:       results,
				UrlCollection: urlCollection,
				Page:          page,
			}

			select {
			case resultChan <- crawlResult:
				if isLocalEnv && atomic.LoadInt32(counter) >= int32(app.engine.DevCrawlLimit) {
					app.Logger.Warn("Dev Crawl limit reached!")
					return
				}
				atomic.AddInt32(counter, 1)
			default:
				app.Logger.Info("Channel is full, dropping Item")
			}
		}
	}
}

type Preference struct {
	MarkAsComplete bool
}

func (app *Crawler) CrawlUrls(collection string, processor interface{}, preferences ...Preference) {
	app.crawlUrlsRecursive(collection, processor, 0, preferences...)
}
func (app *Crawler) crawlUrlsRecursive(collection string, processor interface{}, counter int32, preferences ...Preference) {
	var items []UrlCollection
	var preference Preference
	preference.MarkAsComplete = true
	if len(preferences) > 0 {
		preference = preferences[0]
	}

	urlCollections := app.getUrlCollections(collection)

	var wg sync.WaitGroup
	urlChan := make(chan UrlCollection, len(urlCollections))
	resultChan := make(chan interface{}, len(urlCollections))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, urlCollection := range urlCollections {
		urlChan <- urlCollection
	}
	close(urlChan)

	proxyCount := len(app.engine.ProxyServers)
	batchSize := app.engine.ConcurrentLimit
	totalUrls := len(urlCollections)
	goroutineCount := min(max(proxyCount, 1)*batchSize, totalUrls)

	for i := 0; i < goroutineCount; i++ {
		proxy := Proxy{}
		if proxyCount > 0 {
			proxy = app.engine.ProxyServers[i%proxyCount]
		}
		wg.Add(1)
		go func(proxy Proxy) {
			defer wg.Done()
			app.crawlWorker(ctx, collection, urlChan, resultChan, proxy, processor, isLocalEnv(app.Config.GetString("APP_ENV")), &counter)
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
			case []UrlCollection:
				items = append(items, res...)
				for _, item := range res {
					if item.Parent == "" && collection != baseCollection {
						app.Logger.Fatal("Missing Parent Url, Invalid UrlCollection: %v", item)
						continue
					}
				}
				app.insert(res, v.UrlCollection.Url)
				if preference.MarkAsComplete {
					err := app.markAsComplete(v.UrlCollection.Url, collection)
					if err != nil {
						app.Logger.Error(err.Error())
						continue
					}
				}
				app.Logger.Info("(%d) :%s: Found From [%s => %s]", len(res), app.collection, collection, v.UrlCollection.Url)
			}
		}
	}

	if isLocalEnv(app.Config.GetString("APP_ENV")) && atomic.LoadInt32(&counter) >= int32(app.engine.DevCrawlLimit) {
		cancel()
		return
	}
	if len(urlCollections) > 0 {
		app.crawlUrlsRecursive(collection, processor, counter, preference)
	}
	if len(items) > 0 {
		app.Logger.Info("Total :%s: = (%d)", app.collection, len(items))
	}
}

// CrawlPageDetail initiates the crawling process for detailed page information from the specified collection.
// It distributes the work among multiple goroutines and uses proxies if available.
func (app *Crawler) CrawlPageDetail(collection string, mustRequiredFields ...string) {
	app.CrawlPageDetailRecursive(collection, 0, 0, mustRequiredFields...)
	exportProductDetailsToCSV(app, app.collection, 1)
}

func (app *Crawler) CrawlPageDetailRecursive(collection string, total int, counter int32, mustRequiredFields ...string) {
	urlCollections := app.getUrlCollections(collection)
	var wg sync.WaitGroup
	urlChan := make(chan UrlCollection, len(urlCollections))
	resultChan := make(chan interface{}, len(urlCollections))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, urlCollection := range urlCollections {
		urlChan <- urlCollection
	}
	close(urlChan)

	proxyCount := len(app.engine.ProxyServers)
	batchSize := app.engine.ConcurrentLimit
	totalUrls := len(urlCollections)
	goroutineCount := min(max(proxyCount, 1)*batchSize, totalUrls) // Determine the required number of goroutines

	for i := 0; i < goroutineCount; i++ {
		proxy := Proxy{}
		if proxyCount > 0 {
			proxy = app.engine.ProxyServers[i%proxyCount]
		}
		wg.Add(1)
		go func(proxy Proxy) {
			defer wg.Done()
			app.crawlWorker(ctx, collection, urlChan, resultChan, proxy, app.ProductDetailSelector, isLocalEnv(app.Config.GetString("APP_ENV")), &counter)
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
				invalidFields, unknownFields := validateRequiredFields(res, mustRequiredFields)
				if len(unknownFields) > 0 {
					app.Logger.Error("Unknown fields provided: %v", unknownFields)
					continue
				}
				if len(invalidFields) > 0 {
					app.Logger.Html(v.Page, fmt.Sprintf("Validation failed from URL: %v. Missing value for required fields: %v", v.UrlCollection.Url, invalidFields))
					err := app.markAsError(v.UrlCollection.Url, collection)
					if err != nil {
						app.Logger.Info(err.Error())
						return
					}
					continue
				}

				app.saveProductDetail(res)
				if !isLocalEnv(app.Config.GetString("APP_ENV")) {
					err := app.submitProductData(res)
					if err != nil {
						app.Logger.Fatal("Failed to submit product data to API Server: %v", err)
						err := app.markAsError(v.UrlCollection.Url, collection)
						if err != nil {
							app.Logger.Info(err.Error())
							return
						}
					}
				}

				err := app.markAsComplete(v.UrlCollection.Url, collection)
				if err != nil {
					app.Logger.Error(err.Error())
					continue
				}
				total++
			}
		}
	}
	if isLocalEnv(app.Config.GetString("APP_ENV")) && atomic.LoadInt32(&counter) >= int32(app.engine.DevCrawlLimit) {
		app.Logger.Info("Total %v %v Inserted ", total, app.collection)
		cancel()
		return
	}
	if len(urlCollections) > 0 {
		app.CrawlPageDetailRecursive(collection, total, counter, mustRequiredFields...)
	}
}

// validateRequiredFields checks if the required fields are non-empty in the ProductDetail struct.
// Returns two slices: one for invalid fields and one for unknown fields.
func validateRequiredFields(product *ProductDetail, requiredFields []string) ([]string, []string) {
	var invalidFields []string
	var unknownFields []string

	v := reflect.ValueOf(*product)
	t := v.Type()

	for _, field := range requiredFields {
		f, ok := t.FieldByName(field)
		if !ok {
			unknownFields = append(unknownFields, field)
			continue
		}
		fieldValue := v.FieldByName(field)
		if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
			invalidFields = append(invalidFields, f.Name)
		}
	}
	return invalidFields, unknownFields
}

// PageSelector adds a new URL selector to the crawler.
func (app *Crawler) PageSelector(selector UrlSelector) *Crawler {
	app.UrlSelectors = append(app.UrlSelectors, selector)
	return app
}

// StartUrlCrawling initiates the URL crawling process for all added selectors.
func (app *Crawler) StartUrlCrawling() *Crawler {
	for _, selector := range app.UrlSelectors {
		app.Collection(selector.ToCollection).
			CrawlUrls(selector.FromCollection, selector)
	}
	return app
}
