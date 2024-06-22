package aqua

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "aqua",
		URL:  "https://aqua-has.com",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			BoostCrawling:   true,
			DevCrawlLimit:   300,
			ConcurrentLimit: 10,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
