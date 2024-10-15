package kitamura

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kitamura",
		URL:  "https://shop.kitamura.jp/",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(true),
			DevCrawlLimit:   50,
			StgCrawlLimit:   50,
			ConcurrentLimit: 15,
			Timeout:         30,
			Adapter:         ninjacrawler.String(ninjacrawler.PlayWrightEngine),
			ProxyStrategy:   ninjacrawler.ProxyStrategyRotation,
			ErrorCodes:      []int{403, 429},
			BlockResources:  true,
			BlockedURLs: []string{
				"https://translate.google.com/translate_a/element.js?cb=googleTranslateElementInit",
				"https://shop.kitamura.jp/ec/js/chatbot.js",
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
