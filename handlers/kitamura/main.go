package kitamura

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kitamura",
		URL:  "https://shop.kitamura.jp/",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(true),
			DevCrawlLimit:   999999,
			ConcurrentLimit: 1,
			//WaitForDynamicRendering: true,
			BlockResources: true,
			//ProxyServers: []ninjacrawler.Proxy{
			//	{
			//		Server: "",
			//	},
			//},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler: UrlHandler,
			//ProductHandler: ProductHandler,
		},
	}
}
