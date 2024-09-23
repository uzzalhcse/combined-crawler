package aqua

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "snowpeak",
		URL:  "https://ec.snowpeak.co.jp/snowpeak/ja/%E3%82%AD%E3%83%A3%E3%83%B3%E3%83%97/c/2010000",
		Engine: ninjacrawler.Engine{
			IsDynamic:     ninjacrawler.Bool(true),
			DevCrawlLimit: 300,
			ProxyStrategy: ninjacrawler.ProxyStrategyRotation,
			ProxyServers: []ninjacrawler.Proxy{
				{
					Server:   "http://45.196.54.85:6664",
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
