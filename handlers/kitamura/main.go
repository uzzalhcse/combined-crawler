package kitamura

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kitamura",
		URL:  "https://shop.kitamura.jp/",
		Engine: ninjacrawler.Engine{
			//IsDynamic:       ninjacrawler.Bool(true),
			DevCrawlLimit:   999999,
			ConcurrentLimit: 5,
			Timeout:         90,
			//BlockResources:  true,
			Provider: "zenrows",
			ProviderOption: ninjacrawler.ProviderQueryOption{
				JsRender:       true,
				CustomHeaders:  true,
				OriginalStatus: true,
				//Wait:           5000,
				PremiumProxy: true,
				ProxyCountry: "jp",
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
