package osg

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {

	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categoryHandler,
			Preference:       ninjacrawler.Preference{DoNotMarkAsComplete: true},
		},
	})

}

func categoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	urlCollections := []ninjacrawler.UrlCollection{}
	ctx.Document.Find("nav#search-category-tab ul li a").Each(func(i int, s *goquery.Selection) {

		toolGroup := s.Find("span").Text()
		urStr := fmt.Sprintf("/product-search/catalog/item/?tool-group=%s&abbreviation_method=partial_match&page_no=1", toolGroup)
		urlCollections = append(urlCollections, ninjacrawler.UrlCollection{
			Url:      ctx.App.GetFullUrl(urStr),
			MetaData: nil,
			Parent:   ctx.UrlCollection.Url,
		})
	})
	return urlCollections
}
