package kojima

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kojima",
		URL:  "https://www.kojima.net",
		Engine: ninjacrawler.Engine{
			IsDynamic:          ninjacrawler.Bool(false),
			DevCrawlLimit:      10,
			StgCrawlLimit:      100,
			ConcurrentLimit:    30,
			SleepAfter:         50,
			SleepDuration:      30,
			Timeout:            300,
			RetrySleepDuration: 1,
			ProxyStrategy:      ninjacrawler.ProxyStrategyRotation,
			ErrorCodes:         []int{403, 429},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
