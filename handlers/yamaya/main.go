package yamaya

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "yamaya1",
		URL:  "https://drive.yamaya.jp/index.php",
		Engine: ninjacrawler.Engine{
			IsDynamic: false,
			//BoostCrawling:   true,
			DevCrawlLimit:   200,
			ConcurrentLimit: 2,
			SleepAfter:      10,
			ProxyServers: []ninjacrawler.Proxy{
				{
					Server: "http://35.200.67.230:3000",
				},
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
