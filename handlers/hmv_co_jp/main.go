package hmv_co_jp

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "hmv_co_jp",
		URL:  "https://www.hmv.co.jp",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(true),
			DevCrawlLimit:   0,
			ConcurrentLimit: 50,
			StgCrawlLimit:   0,
			SleepAfter:      300,
			SleepDuration:   30,
			Timeout:         300,
			BlockResources:  true,
			ProxyStrategy:   ninjacrawler.ProxyStrategyRotation,

			Adapter: ninjacrawler.String(ninjacrawler.RodEngine),
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}

/*
https://www.hmv.co.jp/books/
https://www.hmv.co.jp/recordshop
https://www.hmv.co.jp/toy/
https://www.hmv.co.jp/kaitori
https://www.hmv.co.jp/books/genre_%E6%96%87%E8%8A%B8__5_410_0/
https://www.hmv.co.jp/books/genre_%E6%96%87%E8%8A%B8_%E5%9B%BD%E5%86%85%E5%B0%8F%E8%AA%AC_5_410_411/
https://www.hmv.co.jp/select/vinyl/list/sort/recommended/?style=10909&theme=10909001
https://www.hmv.co.jp/books/genre_%E6%96%87%E8%8A%B8_%E6%B5%B7%E5%A4%96%E5%B0%8F%E8%AA%AC_5_410_412/
*/
