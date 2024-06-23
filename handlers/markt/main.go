package markt

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "markt",
		URL:  "https://markt-mall.jp/",
		Engine: ninjacrawler.Engine{
			DevCrawlLimit:   3,
			ConcurrentLimit: 1,
			IsDynamic:       true,
			BlockResources:  true,
			Timeout:         50,
			//BoostCrawling:           true,
			SleepAfter:              900,
			WaitForDynamicRendering: false,
		},
		Handler: ninjacrawler.Handler{
			//UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
