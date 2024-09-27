package panasonic_ec

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "panasonic_ec2",
		URL:  "https://ec-plus.panasonic.jp",
		Engine: ninjacrawler.Engine{
			DevCrawlLimit:   100,
			ConcurrentLimit: 3,
			SleepAfter:      100,
			Timeout:         60,
			BlockResources:  true,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
