package panasonic

import (
	"combined-crawler/pkg/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "panasonic",
		URL:  "https://panasonic.jp/products.html",
		Engine: ninjacrawler.Engine{
			IsDynamic:          ninjacrawler.Bool(false),
			DevCrawlLimit:      200,
			ConcurrentLimit:    1,
			SleepAfter:         10,
			Timeout:            30,  // 30 seconds
			SleepDuration:      120, // 2 minutes
			RetrySleepDuration: 5,
			ProxyStrategy:      ninjacrawler.ProxyStrategyConcurrency,
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
				{
					Server:   "socks5://192.53.66.127:6233",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://45.196.63.200:6834",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://103.130.178.96:5760",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://69.91.142.162:7654",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://5.59.251.210:6249",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://63.141.62.245:6538",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://130.180.233.112:7683",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://62.164.242.146:8723",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://216.98.254.114:6424",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://192.46.189.77:6070",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://156.238.176.48:6730",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://193.160.80.113:6381",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://154.194.24.55:5665",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://192.46.203.195:6161",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://192.46.188.73:5732",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://193.160.83.216:6537",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://208.66.76.204:6128",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://103.196.9.44:6618",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://130.180.232.196:8634",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://69.91.142.35:7527",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "socks5://130.180.235.91:5811",
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
