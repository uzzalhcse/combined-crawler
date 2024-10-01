package ninjacrawler

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func (app *Crawler) Crawl(configs []ProcessorConfig) {
	for _, config := range configs {
		app.Logger.Summary("Starting: %s Crawler", config.OriginCollection)
		app.overrideEngineDefaults(app.engine, &config.Engine)
		app.toggleClient()
		total := int32(0)
		crawlLimit := 0
		if app.isLocalEnv && app.engine.DevCrawlLimit > 0 {
			crawlLimit = app.engine.DevCrawlLimit
		} else if !app.isLocalEnv && app.engine.StgCrawlLimit > 0 {
			crawlLimit = app.engine.StgCrawlLimit
		}
		for {
			productList := app.getUrlCollections(config.OriginCollection)
			if len(productList) == 0 {
				break
			}

			shouldContinue := app.processUrlsWithProxies(productList, config, &total, crawlLimit)

			if !shouldContinue {
				app.Logger.Debug("Crawl limit of %d reached, stopping...", crawlLimit)
				break
			}
		}
		dataCount := app.GetDataCount(config.Entity)
		app.Logger.Summary("[Total (%s) :%s: found from :%s:]", dataCount, config.Entity, config.OriginCollection)

		if errCount := app.GetErrorDataCount(config.OriginCollection); errCount > 0 {
			app.Logger.Summary("Error count: %d", errCount)
		}
		// Consolidate similar cases in switch statement
		switch config.Processor.(type) {
		case ProductDetailSelector, ProductDetailApi, func(CrawlerContext, func([]ProductDetailSelector, string)) error:
			dataCount := app.GetDataCount(config.Entity)
			app.Logger.Summary("Data count: %s", dataCount)
			exportProductDetailsToCSV(app, config.Entity, 1)
		}
	}
}

func (app *Crawler) processUrlsWithProxies(urls []UrlCollection, config ProcessorConfig, total *int32, crawlLimit int) bool {
	var wg sync.WaitGroup
	proxies := app.engine.ProxyServers
	shouldContinue := true
	reqCount := int32(0)
	totalReqCount := 0
	batchCount := 0

	for batchIndex := 0; batchIndex < len(urls); batchIndex += app.engine.ConcurrentLimit {
		batchCount++

		if !shouldContinue {
			break
		}

		proxyIndex := 0
		proxy := Proxy{}
		if len(proxies) > 0 && app.engine.ProxyStrategy == ProxyStrategyConcurrency {
			proxyIndex = totalReqCount % len(proxies)
			proxy = proxies[proxyIndex]
		} else if len(proxies) > 0 && app.engine.ProxyStrategy == ProxyStrategyRotation {
			proxyIndex = int(atomic.LoadInt32(&app.lastWorkingProxyIndex))
			proxy = proxies[proxyIndex]
		}
		app.OpenBrowsers(proxy)

		for i := batchIndex; i < batchIndex+app.engine.ConcurrentLimit && i < len(urls); i++ {
			if crawlLimit > 0 && atomic.LoadInt32(total) >= int32(crawlLimit) {
				shouldContinue = false
				break
			}

			url := urls[i]
			wg.Add(1)

			go func(urlCollection UrlCollection, proxyIndex int) {
				defer func() {
					if r := recover(); r != nil {
						app.HandlePanic(r)
					}
					wg.Done()
				}()
				defer func() {
					app.ClosePages()
				}()

				// Inside goroutine, monitor CPU and RAM usage periodically
				ticker := time.NewTicker(2 * time.Second) // Check system usage every 2 seconds
				defer ticker.Stop()

				done := make(chan struct{})
				go func() {
					for {
						select {
						case <-ticker.C:
							// Check system usage dynamically and take action if necessary
							if app.isCpuUsageHigh() || app.isRamUsageHigh() {
								app.Logger.Warn("CPU or RAM usage exceeds threshold, pausing execution...")
								time.Sleep(5 * time.Second) // Pause for a short time
							}
						case <-done:
							return
						}
					}
				}()

				if crawlLimit > 0 && atomic.AddInt32(total, 1) > int32(crawlLimit) {
					atomic.AddInt32(total, -1)
					shouldContinue = false
					close(done)
					return
				}
				atomic.AddInt32(&reqCount, 1)

				if reqCount > 0 && atomic.LoadInt32(&reqCount)%int32(app.engine.SleepAfter) == 0 {
					app.Logger.Info("Sleeping %d seconds after %d operations", app.engine.SleepDuration, app.engine.SleepAfter)
					time.Sleep(time.Duration(app.engine.SleepDuration) * time.Second)
				}
				app.OpenPages()
				app.crawlWithProxies(urlCollection, config, proxies, proxyIndex, batchCount, 0)

				close(done) // Stop monitoring after the work is done
			}(url, proxyIndex)
		}

		wg.Wait()
		app.CloseBrowsers()
		totalReqCount++
	}

	return shouldContinue
}

func (app *Crawler) crawlWithProxies(urlCollection UrlCollection, config ProcessorConfig, proxies []Proxy, proxyIndex, batchCount, attempt int) {
	proxy := Proxy{}
	if len(proxies) > 0 {
		proxy = proxies[proxyIndex]
	}
	app.CurrentCollection = config.OriginCollection
	app.CurrentUrlCollection = urlCollection
	app.CurrentProxy = proxy
	preHandlerError := false
	if config.Preference.PreHandlers != nil {
		for _, preHandler := range config.Preference.PreHandlers {
			err := preHandler(PreHandlerContext{UrlCollection: urlCollection, App: app})
			if err != nil {
				preHandlerError = true
			}
		}
	}
	if !preHandlerError {
		// Crawl worker execution
		ctx, err := app.handleCrawlWorker(config, proxy, urlCollection)
		if err != nil {
			if strings.Contains(err.Error(), "StatusCode:404") {
				// Mark as max error and stop retrying
				if markMaxErr := app.MarkAsMaxErrorAttempt(urlCollection.Url, config.OriginCollection, err.Error()); markMaxErr != nil {
					app.Logger.Error("markMaxErr: ", markMaxErr.Error())
					return
				}
			} else if strings.Contains(err.Error(), "isRetryable") {
				// Rotate proxy if it's a retryable error
				if len(proxies) > 0 && app.engine.ProxyStrategy == ProxyStrategyRotation {
					nextProxyIndex := (proxyIndex + 1) % len(proxies)
					app.Logger.Summary("Error with proxy %s: %v. Retrying with a different proxy: %s", proxy.Server, err.Error(), proxies[nextProxyIndex].Server)

					// Check if we have exhausted all proxies
					if attempt >= len(proxies) {
						app.Logger.Info("All proxies exhausted.")
						return
						//time.Sleep(1 * time.Hour)
						//app.crawlWithProxies(urlCollection, config, proxies, 0, batchCount, 0) // Restart with the first proxy
					} else {
						app.crawlWithProxies(urlCollection, config, proxies, nextProxyIndex, batchCount, attempt+1) // Retry with the next proxy
					}
					return
				}
				if app.engine.RetrySleepDuration > 0 {
					app.HandleThrottling(1, urlCollection.StatusCode)
				}
				if markErr := app.MarkAsError(urlCollection.Url, config.OriginCollection, err.Error()); markErr != nil {
					app.Logger.Error("markErr: ", markErr.Error())
					return
				}
			} else {
				if markErr := app.MarkAsError(urlCollection.Url, config.OriginCollection, err.Error()); markErr != nil {
					app.Logger.Error("markErr: ", markErr.Error())
				}
			}
			app.Logger.Error("Error crawling %s: %v", urlCollection.Url, err)
			return
		}
		// Process successful crawl
		app.extract(config, *ctx)
		// Update last working proxy index on success
		atomic.StoreInt32(&app.lastWorkingProxyIndex, int32(proxyIndex))

	}
}
