package ninjacrawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"strings"
	"sync/atomic"
	"time"
)

func (app *Crawler) crawlWorker(ctx context.Context, processorConfig ProcessorConfig, urlChan <-chan UrlCollection, resultChan chan<- interface{}, isLocalEnv bool, counter *int32, currentProxyIndex int) {
	var page playwright.Page
	var browser playwright.Browser
	var err error
	var doc *goquery.Document
	var apiResponse map[string]interface{}

	// Used to track the proxy index
	proxyIndex := 0
	if app.engine.ProxyStrategy == ProxyStrategyConcurrency && currentProxyIndex > 0 {
		proxyIndex = currentProxyIndex
	}
	usedProxies := make(map[int]bool)

	rotateProxy := func() Proxy {
		proxyIndex = (proxyIndex + 1) % len(app.engine.ProxyServers)
		usedProxies[proxyIndex] = true
		app.Logger.Warn(fmt.Sprintf("Rotating proxy proxyIndex %d", proxyIndex))
		return app.engine.ProxyServers[proxyIndex]
	}

	// Get the initial proxy
	currentProxy := app.engine.ProxyServers[proxyIndex]
	usedProxies[proxyIndex] = true
	app.CurrentProxy = currentProxy

	if *app.engine.IsDynamic {
		browser, page, err = app.GetBrowserPage(app.pw, app.engine.BrowserType, currentProxy)
		if err != nil {
			app.Logger.Fatal(err.Error())
		}
		defer browser.Close()
		defer page.Close()
	}

	operationCount := 0 // Initialize operation count

	for {
		select {
		case <-ctx.Done(): // Handle context timeout or cancellation
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				app.Logger.Warn("Crawl worker timed out")
			} else if errors.Is(ctx.Err(), context.Canceled) {
				app.Logger.Info("Crawl worker was canceled")
			}
			return
		case urlCollection, more := <-urlChan:
			if !more {
				return
			}
			if app.engine.ProxyStrategy == ProxyStrategyRotation {
				if urlCollection.StatusCode == 403 {
					if app.engine.RetrySleepDuration > 0 && app.engine.Provider == "zenrows" && urlCollection.StatusCode >= 400 && urlCollection.StatusCode < 500 && urlCollection.StatusCode != 404 {
						app.HandleThrottling(urlCollection.Attempts, urlCollection.StatusCode)
					}
					//app.HandleThrottling(urlCollection.Attempts, urlCollection.StatusCode)
					// Rotate the proxy on receiving a 403
					currentProxy = rotateProxy()
					app.CurrentProxy = currentProxy
					app.Logger.Debug("Received 403. Rotating proxy to %s", currentProxy.Server)

					if *app.engine.IsDynamic {
						browser, page, err = app.GetBrowserPage(app.pw, app.engine.BrowserType, currentProxy)
						if err != nil {
							app.Logger.Fatal(err.Error())
						}
					}
				}
			}
			preHandlerError := false
			if processorConfig.Preference.PreHandlers != nil { // Execute pre handlers
				for _, preHandler := range processorConfig.Preference.PreHandlers {
					err := preHandler(PreHandlerContext{UrlCollection: urlCollection, App: app})
					if err != nil {
						preHandlerError = true
					}
				}
			}
			if preHandlerError {
				continue
			}
			crawlableUrl := urlCollection.Url
			if urlCollection.ApiUrl != "" {
				crawlableUrl = urlCollection.ApiUrl
			}
			if urlCollection.CurrentPageUrl != "" {
				crawlableUrl = urlCollection.CurrentPageUrl
			}

			if currentProxy.Server != "" {
				app.Logger.Info("Crawling :%s: %s using Proxy %s", processorConfig.OriginCollection, crawlableUrl, currentProxy.Server)
			} else {
				app.Logger.Info("Crawling :%s: %s", processorConfig.OriginCollection, crawlableUrl)
			}
			if *app.engine.IsDynamic {
				doc, err = app.NavigateToURL(page, crawlableUrl)
			} else {
				switch processorConfig.Processor.(type) {
				case ProductDetailApi:
					apiResponse, err = app.NavigateToApiURL(app.httpClient, crawlableUrl, currentProxy)
				default:
					doc, err = app.NavigateToStaticURL(app.httpClient, crawlableUrl, currentProxy)
				}
			}

			if err != nil {
				if strings.Contains(err.Error(), "StatusCode:404") {
					if markMaxErr := app.MarkAsMaxErrorAttempt(urlCollection.Url, processorConfig.OriginCollection, err.Error()); markMaxErr != nil {
						app.Logger.Error("markMaxErr: ", markMaxErr.Error())
						return
					}
				} else {
					if markErr := app.MarkAsError(urlCollection.Url, processorConfig.OriginCollection, err.Error()); markErr != nil {
						app.Logger.Error("markErr: ", markErr.Error())
						return
					}
				}
				app.Logger.Error("Error crawling %s: %v", urlCollection.Url, err)
				continue
			}

			crawlerCtx := CrawlerContext{
				App:           app,
				Document:      doc,
				UrlCollection: urlCollection,
				Page:          page,
				ApiResponse:   apiResponse,
			}

			var results interface{}
			switch v := processorConfig.Processor.(type) {
			case func(CrawlerContext) []UrlCollection:
				var collections []UrlCollection
				collections = v(crawlerCtx)

				for _, item := range collections {
					if item.Parent == "" && processorConfig.OriginCollection != baseCollection {
						app.Logger.Fatal("Missing Parent Url, Invalid OriginCollection: %v", item)
						continue
					}
				}
				app.insert(processorConfig.Entity, collections, urlCollection.Url)
				if !processorConfig.Preference.DoNotMarkAsComplete {
					err := app.markAsComplete(urlCollection.Url, processorConfig.OriginCollection)
					if err != nil {
						app.Logger.Error(err.Error())
						continue
					}
				}
				atomic.AddInt32(counter, 1)
			case func(CrawlerContext, func([]UrlCollection, string)) error:
				shouldMarkAsComplete := true
				handleErr := v(crawlerCtx, func(collections []UrlCollection, currentPageUrl string) {
					for _, item := range collections {
						if item.Parent == "" && processorConfig.OriginCollection != baseCollection {
							app.Logger.Fatal("Missing Parent Url, Invalid OriginCollection: %v", item)
							continue
						}
					}
					if currentPageUrl != "" && currentPageUrl != urlCollection.Url {
						shouldMarkAsComplete = false
						currentPageErr := app.SyncCurrentPageUrl(urlCollection.Url, currentPageUrl, processorConfig.OriginCollection)
						if currentPageErr != nil {
							app.Logger.Fatal(currentPageErr.Error())
							return
						}
					} else {
						shouldMarkAsComplete = true
						atomic.AddInt32(counter, 1)
					}
					app.insert(processorConfig.Entity, collections, urlCollection.Url)
				})
				if handleErr != nil {
					markAsError := app.MarkAsError(urlCollection.Url, processorConfig.OriginCollection, handleErr.Error())
					if markAsError != nil {
						app.Logger.Info(markAsError.Error())
						return
					}
					app.Logger.Error(handleErr.Error())
				} else {
					if !processorConfig.Preference.DoNotMarkAsComplete && shouldMarkAsComplete {
						err := app.markAsComplete(urlCollection.Url, processorConfig.OriginCollection)
						if err != nil {
							app.Logger.Error(err.Error())
							continue
						}
					}
				}

			case UrlSelector:
				var collections []UrlCollection
				collections = app.processDocument(doc, v, urlCollection)

				for _, item := range collections {
					if item.Parent == "" && processorConfig.OriginCollection != baseCollection {
						app.Logger.Fatal("Missing Parent Url, Invalid OriginCollection: %v", item)
						continue
					}
				}
				app.insert(processorConfig.Entity, collections, urlCollection.Url)

				if !processorConfig.Preference.DoNotMarkAsComplete {
					err := app.markAsComplete(urlCollection.Url, processorConfig.OriginCollection)
					if err != nil {
						app.Logger.Error(err.Error())
						continue
					}
				}
				atomic.AddInt32(counter, 1)

			case func(CrawlerContext, func([]ProductDetailSelector, string)) error:
				shouldMarkAsComplete := true
				handleErr := v(crawlerCtx, func(collections []ProductDetailSelector, currentPageUrl string) {
					if currentPageUrl != "" && currentPageUrl != urlCollection.Url {
						shouldMarkAsComplete = false
						currentPageErr := app.SyncCurrentPageUrl(urlCollection.Url, currentPageUrl, processorConfig.OriginCollection)
						if currentPageErr != nil {
							app.Logger.Fatal(currentPageErr.Error())
							return
						}
					} else {
						shouldMarkAsComplete = true
						atomic.AddInt32(counter, 1)
					}

					for _, collection := range collections {
						res := crawlerCtx.handleProductDetail(collection)
						result := CrawlResult{
							Results:       res,
							UrlCollection: urlCollection,
							Page:          page,
							Document:      doc,
						}
						err := app.handleProductDetail(res, processorConfig, result)
						if err != nil {
							app.Logger.Error(err.Error())
							continue
						}
					}
				})
				if handleErr != nil {
					markAsError := app.MarkAsError(urlCollection.Url, processorConfig.OriginCollection, handleErr.Error())
					if markAsError != nil {
						app.Logger.Info(markAsError.Error())
						return
					}
					app.Logger.Error(handleErr.Error())
				} else {
					if !processorConfig.Preference.DoNotMarkAsComplete && shouldMarkAsComplete {
						err := app.markAsComplete(urlCollection.Url, processorConfig.OriginCollection)
						if err != nil {
							app.Logger.Error(err.Error())
							continue
						}
					}
				}
			case ProductDetailSelector:
				results = crawlerCtx.handleProductDetail(processorConfig.Processor)
			case ProductDetailApi:
				results = crawlerCtx.handleProductDetailApi(processorConfig.Processor)

			default:
				app.Logger.Fatal("Unsupported processor type: %T", processorConfig.Processor)
			}

			crawlResult := CrawlResult{
				Results:       results,
				UrlCollection: urlCollection,
				Page:          page,
				Document:      doc,
			}

			select {
			case resultChan <- crawlResult:
				if isLocalEnv && atomic.LoadInt32(counter) >= int32(app.engine.DevCrawlLimit) {
					//app.Logger.Warn("Dev Crawl limit %d reached!...", atomic.LoadInt32(counter))
					return
				}
				atomic.AddInt32(counter, 1)
			default:
				app.Logger.Info("Channel is full, dropping Item")
			}
			if isLocalEnv && atomic.LoadInt32(counter) >= int32(app.engine.DevCrawlLimit) {
				//app.Logger.Warn("Dev Crawl limit %d reached!", atomic.LoadInt32(counter))
				return
			}

			operationCount++                               // Increment the operation count
			if operationCount%app.engine.SleepAfter == 0 { // Apply sleep after a certain number of operations
				app.Logger.Info("Sleeping %d seconds after %d operations", app.engine.SleepDuration, app.engine.SleepAfter)
				time.Sleep(time.Duration(app.engine.SleepDuration) * time.Second)
			}
		}
	}
}
