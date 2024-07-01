package yamaya

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "yamaya",
		URL:  "https://drive.yamaya.jp/index.php",
		Engine: ninjacrawler.Engine{
			IsDynamic: false,
			//BoostCrawling:   true,
			DevCrawlLimit:   500,
			ConcurrentLimit: 1,
			SleepAfter:      5,
			//ProxyServers: []ninjacrawler.Proxy{
			//	{
			//		Server: "http://35.221.68.83:3000",
			//	},
			//},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
