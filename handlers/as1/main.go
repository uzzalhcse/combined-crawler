package as1

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "as1",
		URL:  "https://axel.as-1.co.jp/",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			DevCrawlLimit:   300,
			ConcurrentLimit: 10,
			BlockResources:  true,
			Provider:        "zenrows",
			Timeout:         300, // 5 minute
			ProviderOption: ninjacrawler.ProviderQueryOption{
				JsRender:       true,
				CustomHeaders:  true,
				OriginalStatus: true,
				//WaitFor:        ".af-price > span",
			},
		},
		Handler: ninjacrawler.Handler{
			//UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
