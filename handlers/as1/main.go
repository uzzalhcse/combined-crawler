package as1

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "as2",
		URL:  "https://axel.as-1.co.jp/",
		Engine: ninjacrawler.Engine{
			BrowserType:     "firefox",
			IsDynamic:       true,
			DevCrawlLimit:   300,
			ConcurrentLimit: 1,
			//BoostCrawling:   true,
			BlockResources: true,
			//ProxyServers: []ninjacrawler.Proxy{
			//	{
			//		Server: "http://35.243.121.253:3000",
			//	},
			//},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler: UrlHandler,
			//ProductHandler: ProductHandler,
		},
	}
}
