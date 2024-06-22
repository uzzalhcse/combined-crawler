package kyocera

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kyocera",
		URL:  "https://www.kyocera.co.jp/prdct/tool/category/product",
		Engine: ninjacrawler.Engine{
			BoostCrawling:   true,
			IsDynamic:       false,
			DevCrawlLimit:   300,
			ConcurrentLimit: 15,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
