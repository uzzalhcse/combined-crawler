package yamaya

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "yamaya1",
		URL:  "https://drive.yamaya.jp/index.php",
		Engine: ninjacrawler.Engine{
			IsDynamic: false,
			//BoostCrawling:   true,
			DevCrawlLimit:   100,
			ConcurrentLimit: 4,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler: UrlHandler,
			//ProductHandler: ProductHandler,
		},
	}
}
