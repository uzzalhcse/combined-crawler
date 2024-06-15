package markt

import (
	"github.com/lazuli-inc/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "markt",
		URL:  "https://markt-mall.jp/",
		Engine: ninjacrawler.Engine{
			IsDynamic:       true,
			BoostCrawling:   true,
			ConcurrentLimit: 5,
			DevCrawlLimit:   1,
			BlockResources:  true,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
