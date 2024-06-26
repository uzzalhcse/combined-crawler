package kyocera

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kyocera",
		URL:  "https://www.kyocera.co.jp/prdct/tool/category/product",
		Engine: ninjacrawler.Engine{
			BoostCrawling:           true,
			IsDynamic:               false,
			DevCrawlLimit:           20,
			ConcurrentLimit:         2,
			SleepAfter:              30,
			WaitForDynamicRendering: true,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
