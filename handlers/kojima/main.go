package kojima

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kojima",
		URL:  "https://www.kojima.net",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(false),
			DevCrawlLimit:   10,
			ConcurrentLimit: 5,
			SleepAfter:      500,
			//Provider:        "zenrows",
			SleepDuration: 30,
			//ProviderOption: ninjacrawler.ProviderQueryOption{
			//	UsePremiumProxyRetry: true,
			//	CustomHeaders:        true,
			//	OriginalStatus:       true,
			//},
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
