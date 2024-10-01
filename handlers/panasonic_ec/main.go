package panasonic_ec

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "panasonic_ec",
		URL:  "https://ec-plus.panasonic.jp",
		Engine: ninjacrawler.Engine{
			DevCrawlLimit:   100,
			ConcurrentLimit: 20,
			StgCrawlLimit:   150,
			SleepAfter:      30,
			SleepDuration:   60,
			Timeout:         60,
			//BlockResources:  true,
		},
		Handler: ninjacrawler.Handler{
			//UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
