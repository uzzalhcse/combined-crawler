package kojima

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "kojima",
		URL:  "https://www.kojima.net",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(false),
			DevCrawlLimit:   10,
			ConcurrentLimit: 3,
			SleepAfter:      500,
			//Provider:        "zenrows",
			SleepDuration: 30,
			//ProviderOption: ninjacrawler.ProviderQueryOption{
			//	UsePremiumProxyRetry: true,
			//	CustomHeaders:        true,
			//	OriginalStatus:       true,
			//},
			Timeout:            300,
			RetrySleepDuration: 1,
			ProxyStrategy:      ninjacrawler.ProxyStrategyRotation,
			ProxyServers: []ninjacrawler.Proxy{
				{
					Server:   "http://62.164.231.7:9319",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "http://192.46.190.170:6763",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "http://130.180.233.112:7683",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "http://72.46.139.119:6679",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "http://45.196.43.235:5962",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "http://192.145.71.5:6642",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "http://130.180.231.111:8253",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
				{
					Server:   "http://192.53.66.127:6233",
					Username: "lnvmpyru",
					Password: "5un1tb1azapa",
				},
			},
			ErrorCodes: []int{403, 429},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
