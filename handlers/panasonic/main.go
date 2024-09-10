package panasonic

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "panasonic",
		URL:  "https://panasonic.jp/products.html",
		Engine: ninjacrawler.Engine{
			IsDynamic:          ninjacrawler.Bool(false),
			DevCrawlLimit:      200,
			ConcurrentLimit:    10,
			SleepAfter:         5,
			Timeout:            30,
			RetrySleepDuration: 30,
			Provider:           "zenrows",
			ProviderOption: ninjacrawler.ProviderQueryOption{
				OriginalStatus: true,
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
