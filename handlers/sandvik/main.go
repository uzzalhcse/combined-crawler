package sandvik

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "sandvik",
		URL:  "https://www.sandvik.coromant.com/ja-jp/tools",
		Engine: ninjacrawler.Engine{
			BoostCrawling:  true,
			BlockResources: true,
			//DevCrawlLimit:  100,
			IsDynamic: true,
			CookieConsent: &ninjacrawler.CookieAction{
				ButtonText:                  "Accept Cookies",
				MustHaveSelectorAfterAction: "body .column.grid_12.col-12",
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
