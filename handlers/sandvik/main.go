package sandvik

import (
	"github.com/lazuli-inc/ninjacrawler"
)

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "sandvik",
		URL:  "https://www.sandvik.coromant.com/ja-jp/tools",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			BoostCrawling:   true,
			ConcurrentLimit: 20,
			DevCrawlLimit:   1,
			CookieConsent: &ninjacrawler.CookieAction{
				ButtonText:       "Accept Cookies",
				SleepAfterAction: 7,
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
