package sumitool

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {

	selector := ninjacrawler.UrlSelector{
		Selector:     "ul.category-index-list li.category-index-item",
		FindSelector: "a",
		Attr:         "href",
	}
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        selector,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        selector,
		},
	})

}
