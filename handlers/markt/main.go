package markt

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "markt",
		URL:  "https://markt-mall.jp/",
		Engine: ninjacrawler.Engine{
			DevCrawlLimit:           100,
			ConcurrentLimit:         10,
			IsDynamic:               true,
			BlockResources:          true,
			Timeout:                 50,
			BoostCrawling:           true,
			SleepAfter:              30,
			WaitForDynamicRendering: true,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
