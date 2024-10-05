package panasonic

import (
	"combined-crawler/pkg/ninjacrawler"
	"time"
)

const maxAttempts = 2
const retryDelay = 2 * time.Second

func Crawler() ninjacrawler.CrawlerConfig {
	return ninjacrawler.CrawlerConfig{
		Name: "panasonic",
		URL:  "https://panasonic.jp/products.html",
		Engine: ninjacrawler.Engine{
			IsDynamic:               ninjacrawler.Bool(false),
			ConcurrentLimit:         50,
			SleepAfter:              60,
			SleepDuration:           60,
			StgCrawlLimit:           1000,
			StoreHtml:               ninjacrawler.Bool(false),
			SendHtmlToBigquery:      ninjacrawler.Bool(false),
			ErrorCodes:              []int{403, 429},
			IgnoreRetryOnValidation: ninjacrawler.Bool(true),
			ProxyStrategy:           ninjacrawler.ProxyStrategyRotation,
		},
		Handler: ninjacrawler.Handler{
			UrlHandler:     UrlHandler,
			ProductHandler: ProductHandler,
		},
	}
}
