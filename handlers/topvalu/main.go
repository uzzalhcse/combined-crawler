package topvalu

import (
	"combined-crawler/pkg/ninjacrawler"
)

const (
	name = "topvalu"
	url  = "https://www.topvalu.net"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: name,
		URL:  url,
		Engine: ninjacrawler.Engine{
			BoostCrawling:  false,
			BlockResources: false,
			IsDynamic:      ninjacrawler.Bool(true),
			DevCrawlLimit:  1000000000,
			CookieConsent: &ninjacrawler.CookieAction{
				ButtonText:                  "Accept Cookies",
				MustHaveSelectorAfterAction: "body .header__gnav.gnav",
			},
			//ProxyServers: []ninjacrawler.Proxy{
			//	{
			//		Server:   "http://5.59.251.78:6117",
			//		Username: "lnvmpyru",
			//		Password: "5un1tb1azapa",
			//	},
			//},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductDetailsHandler,
		},
	}
}
