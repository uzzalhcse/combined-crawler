package sandvik

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "sandvik",
		URL:  "https://www.sandvik.coromant.com/ja-jp/tools",
		Engine: ninjacrawler.Engine{
			IsDynamic:       false,
			DevCrawlLimit:   1,
			ConcurrentLimit: 1,
			CookieConsent: &ninjacrawler.CookieAction{
				ButtonText:                  "Accept Cookies",
				MustHaveSelectorAfterAction: "body .page-container",
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
