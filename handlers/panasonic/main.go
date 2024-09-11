package panasonic

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "panasonic",
		URL:  "https://panasonic.jp/products.html",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(false),
			DevCrawlLimit:   200,
			ConcurrentLimit: 3,
			SleepAfter:      5,
			Timeout:         30,  // 30 seconds
			SleepDuration:   120, // 2 minutes
			//RetrySleepDuration: 30,
			ProxyServers: []ninjacrawler.Proxy{
				{
					Server:   "socks5://5.59.251.78:6117",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://5.59.251.19:6058",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://62.164.231.7:9319",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://192.46.190.170:6763",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://130.180.233.112:7683",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://208.73.42.138:9149",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://72.46.139.119:6679",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://45.196.43.235:5962",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://192.145.71.5:6642",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://130.180.231.111:8253",
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
