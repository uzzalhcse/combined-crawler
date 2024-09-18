package topvalu

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        CategoryHandler,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        ProductHandler,
		},
	})
}
