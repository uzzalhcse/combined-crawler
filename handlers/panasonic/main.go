package panasonic

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "panasonic",
		URL:  "https://panasonic.jp/products.html",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			DevCrawlLimit:   200,
			ConcurrentLimit: 10,
			SleepAfter:      10,
			Timeout:         30,
			BoostCrawling:   true,
		},
		Handler: ninjacrawler.Handler{
			//UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
