package kitamura

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kitamura",
		URL:  "https://shop.kitamura.jp/",
		Engine: ninjacrawler.Engine{
			BrowserType:     "webkit",
			IsDynamic:       ninjacrawler.Bool(true),
			DevCrawlLimit:   999999,
			ConcurrentLimit: 1,
			//WaitForDynamicRendering: true,
			WaitForSelector: ninjacrawler.String(".category-item"),
			BlockResources:  true,
			//Provider:       "zenrows",
			//ProviderOption: ninjacrawler.ProviderQueryOption{
			//	JsRender:       true,
			//	CustomHeaders:  true,
			//	OriginalStatus: true,
			//	Wait:           5000,
			//},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler: UrlHandler,
			//ProductHandler: ProductHandler,
		},
	}
}
