package suntory

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "suntory",
		URL:  "https://products.suntory.co.jp?ke=hd",
		Engine: ninjacrawler.Engine{
			IsDynamic:               ninjacrawler.Bool(true),
			DevCrawlLimit:           999999,
			ConcurrentLimit:         5,
			WaitForDynamicRendering: true,
			BlockResources:          true,
			BrowserType:             "firefox",
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
