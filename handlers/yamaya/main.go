package yamaya

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "yamaya",
		URL:  "https://drive.yamaya.jp/index.php",
		Engine: ninjacrawler.Engine{
			IsDynamic: false,
			//BoostCrawling:   true,
			DevCrawlLimit:   10,
			ConcurrentLimit: 5,
			SleepAfter:      50,
			Timeout:         60,
			ProxyServers: []ninjacrawler.Proxy{
				{
					Server: "http://34.48.154.203:3000",
				},
				{
					Server: "http://34.48.157.202:3000",
				},
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
