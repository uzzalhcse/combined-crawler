package as1

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "as2",
		URL:  "https://axel.as-1.co.jp/",
		Engine: ninjacrawler.Engine{
			IsDynamic:       true,
			DevCrawlLimit:   300,
			ConcurrentLimit: 1,
			BoostCrawling:   false,
			BlockResources:  true,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler: UrlHandler,
			//ProductHandler: ProductHandler,
		},
	}
}
