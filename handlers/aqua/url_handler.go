package aqua

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {

	crawler.Logger.Debug("This custom Debug Log for testing")
	categorySelector := ninjacrawler.UrlSelector{
		Selector:     "ul li.category",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
	}
	//productSelector := ninjacrawler.UrlSelector{
	//	Selector:     "div.CategoryTop_Series_Item_Content_List",
	//	SingleResult: false,
	//	FindSelector: "a",
	//	Attr:         "href",
	//}

	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categorySelector,
		},
		//{
		//	Entity:           constant.Products,
		//	OriginCollection: constant.Categories,
		//	Processor:        productSelector,
		//},
	})

}
