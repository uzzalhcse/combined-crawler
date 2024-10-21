package hmv_co_jp

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	categorySelector := ninjacrawler.UrlSelector{
		Selector:     "ul.listSubInnerList li.listSmallSub",
		FindSelector: "a",
		Attr:         "href",
	}
	productSelector := ninjacrawler.UrlSelector{
		Selector:     "ul.resultList li .thumbnailBlock",
		FindSelector: "a",
		Attr:         "href",
	}
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categorySelector,
			Engine:           ninjacrawler.Engine{},
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        productSelector,
			Engine:           ninjacrawler.Engine{},
		},
	})

}
