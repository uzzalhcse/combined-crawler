package panasonic_ec

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "panasonic_ec2",
		URL:  "https://ec-plus.panasonic.jp",
		Engine: ninjacrawler.Engine{
			DevCrawlLimit:   500,
			ConcurrentLimit: 1,
			SleepAfter:      100,
			Timeout:         30,
			BlockResources:  true,
		},
		Handler: ninjacrawler.Handler{
			//UrlHandler: UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
