package asics

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: Name,
		URL:  Url,
		Engine: ninjacrawler.Engine{
			//IsDynamic:          ninjacrawler.Bool(false),
			//DevCrawlLimit:   1000000,
			//Timeout:         60,
			ConcurrentLimit: 1,
			Provider:        "zenrows",
			ProviderOption: ninjacrawler.ProviderQueryOption{
				//JsRender:       false,
				//CustomHeaders:  true,
				//OriginalStatus: true,
				Wait: 2000,
				//PremiumProxy: true,
				//ProxyCountry: "jp",
			},
			//ErrorCodes: []int{403, 429, 422},
			//MaxRetryAttempts: 5,
			//ProxyStrategy: ninjacrawler.ProxyStrategyRotation,
			//ProxyServers:  generateProxy(),
		},
		Handler: ninjacrawler.Handler{
			//UrlHandler:     UrlHandler,
			ProductHandler: ProductDetailsHandler,
		},
	}
}
