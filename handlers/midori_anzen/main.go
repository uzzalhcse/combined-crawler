package midori_anzen

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "midori_anzen",
		URL:  "https://ec.midori-anzen.com/shop/category/category.aspx?plus=0",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			DevCrawlLimit:   300,
			ConcurrentLimit: 10,
			BoostCrawling:   false,
			BlockResources:  true,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
