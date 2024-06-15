package kyocera

import "github.com/lazuli-inc/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kyocera",
		URL:  "https://www.kyocera.co.jp/prdct/tool/category/product",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			BoostCrawling:   true,
			ConcurrentLimit: 20,
			DevCrawlLimit:   2,
			BlockResources:  true,
			BlockedURLs:     []string{"syncsearch.jp"},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
