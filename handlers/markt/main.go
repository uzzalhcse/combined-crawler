package markt

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "markt",
		URL:  "https://markt-mall.jp/",
		Engine: ninjacrawler.Engine{
			DevCrawlLimit:   1,
			ConcurrentLimit: 5,
			IsDynamic:       true,
			BlockResources:  true,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
