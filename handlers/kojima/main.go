package kojima

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kojima",
		URL:  "https://www.kojima.net",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			DevCrawlLimit:   50,
			ConcurrentLimit: 3,
			BlockResources:  true,
			SleepAfter:      20,
			//Provider:        "zenrows",
		},
		Handler: ninjacrawler.Handler{
			UrlHandler: UrlHandler,
			//ProductHandler: ProductHandler,
		},
	}
}
