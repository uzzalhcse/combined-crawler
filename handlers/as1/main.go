package as1

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "as1",
		URL:  "https://axel.as-1.co.jp/",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(true),
			DevCrawlLimit:   300,
			ConcurrentLimit: 15,
			BlockResources:  true,
			Timeout:         300, // 5 minute
			SleepAfter:      50,
			//Provider:        "zenrows",
			//ProviderOption: ninjacrawler.ProviderQueryOption{
			//	JsRender:       true,
			//	CustomHeaders:  true,
			//	OriginalStatus: true,
			//	//WaitFor:        ".af-price > span",
			//},
			Adapter:       ninjacrawler.String(ninjacrawler.PlayWrightEngine),
			ErrorCodes:    []int{403, 429},
			ProxyStrategy: ninjacrawler.ProxyStrategyRotation,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
