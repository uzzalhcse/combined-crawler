package aqua

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "snowpeak",
		URL:  "https://ec.snowpeak.co.jp/snowpeak/ja/%E3%82%AD%E3%83%A3%E3%83%B3%E3%83%97/c/2010000",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(true),
			DevCrawlLimit:   300,
			ConcurrentLimit: 1,
			ProxyStrategy:   ninjacrawler.ProxyStrategyRotation,
			ProxyServers: []ninjacrawler.Proxy{
				{
					Server:   "http://5.59.251.78:6117",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
