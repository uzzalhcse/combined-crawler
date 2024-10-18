package sony

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "sony",
		URL:  "https://www.sony.jp/products_menu.html?s_pid=jp_top_PRODUCTS_ICHIRAN",
		Engine: ninjacrawler.Engine{
			//IsDynamic:               false,
			DevCrawlLimit:           2000,
			ConcurrentLimit:         2,
			SleepAfter:              30,
			Timeout:                 60,
			WaitForDynamicRendering: true,
			BlockResources:          true,
			BlockedURLs: []string{
				"https://mboxedge38.tt.omtrdc.net",
				"https://play.google.com",
				"https://www.youtube.com",
				"https://www.sony.jp/script",
			},
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
