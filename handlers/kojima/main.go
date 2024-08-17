package kojima

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kojima_v2",
		URL:  "https://www.kojima.net",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			DevCrawlLimit:   50,
			ConcurrentLimit: 3,
			BlockResources:  true,
			SleepAfter:      10,
			Provider:        "zenrows",
			ProviderOption:  ninjacrawler.ProviderQueryOption{UsePremiumProxyRetry: true},
			Timeout:         300,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
