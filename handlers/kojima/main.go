package kojima

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kojima",
		URL:  "https://www.kojima.net",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(false),
			DevCrawlLimit:   50,
			ConcurrentLimit: 3,
			SleepAfter:      10,
			Provider:        "zenrows",
			SleepDuration:   30,
			ProviderOption: ninjacrawler.ProviderQueryOption{
				UsePremiumProxyRetry: true,
				CustomHeaders:        true,
				OriginalStatus:       true,
			},
			Timeout:            300,
			RetrySleepDuration: 30,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
