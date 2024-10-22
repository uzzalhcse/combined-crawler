package asics

import "combined-crawler/pkg/ninjacrawler"

func UrlHandler(crawler *ninjacrawler.Crawler) {
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        CategoryHandler,
		},
		{
			Entity:           Products,
			OriginCollection: Categories,
			Processor:        ProductUrlHandler,
		},
	})
}
