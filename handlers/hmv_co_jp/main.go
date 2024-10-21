package hmv_co_jp

import "combined-crawler/pkg/ninjacrawler"

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "hmv_co_jp",
		URL:  "https://www.hmv.co.jp",
		Engine: ninjacrawler.Engine{
			IsDynamic:       ninjacrawler.Bool(true),
			DevCrawlLimit:   0,
			ConcurrentLimit: 2,
			StgCrawlLimit:   0,
			SleepAfter:      50,
			SleepDuration:   30,
			Timeout:         600,
			//BlockResources:  true,
			ProxyStrategy: ninjacrawler.ProxyStrategyRotation,

			//Adapter: ninjacrawler.String(ninjacrawler.RodEngine),
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}

/*
https://books.rakuten.co.jp/genrelist/e-book.html?l-id=header-subnavi-ebook-genrelist
https://books.rakuten.co.jp/ranking/hourly/101/?l-id=header-subnavi-ebook-ranking#!/
https://books.rakuten.co.jp/book/sheet-of-music/?l-id=header-subnavi-book-g001018
https://books.rakuten.co.jp/book/author/?l-id=header-subnavi-book-author
https://books.rakuten.co.jp/ranking/hourly/006/?l-id=header-subnavi-game-ranking#!/
https://books.rakuten.co.jp/download/?l-id=header-subnavi-software-download
https://books.rakuten.co.jp/info/special-price-sale/book/?l-id=header-subnavi-book-special-price-sale

https://books.rakuten.co.jp/event/limited-item/?l-id=header-subnavi-book-limited-item
https://books.rakuten.co.jp/event/book/?l-id=header-subnavi-book-campaign
https://books.rakuten.co.jp/event/limited-item/?l-id=header-subnavi-book-limited-item
https://books.rakuten.co.jp/event/e-book/camp-bestprice/?l-id=header-subnavi-ebook-bestprice
https://books.rakuten.co.jp/event/e-book/free/?l-id=header-subnavi-ebook-free
https://books.rakuten.co.jp/event/e-book/ereaders/?l-id=header-subnavi-ebook-ereaders
https://books.rakuten.co.jp/event/magazine/?l-id=header-subnavi-magazine-campaign



need to exclude some coupons pages like
https://books.rakuten.co.jp/event/coupon/?shop=kobo&l-id=header-subnavi-ebook-coupon

found  some unknown domain so we need to filter those domains except rakuten

count the total data
https://books.rakuten.co.jp/json/genre/000
https://books.rakuten.co.jp/json/genre/001
https://books.rakuten.co.jp/json/genre/002
https://books.rakuten.co.jp/json/genre/003
https://books.rakuten.co.jp/json/genre/004
https://books.rakuten.co.jp/json/genre/005
https://books.rakuten.co.jp/json/genre/006
https://books.rakuten.co.jp/json/genre/007
https://books.rakuten.co.jp/json/genre/101
*/
