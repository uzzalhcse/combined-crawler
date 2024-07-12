package osg

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {

	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categoryHandler,
		},
	})

}

func categoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	urlCollections := []ninjacrawler.UrlCollection{}
	ctx.Document.Find("div[id^=search-category-criteria-]").Each(func(i int, s *goquery.Selection) {
		// Extract the tool-group value from the hidden input
		toolGroup, exists := s.Find("input[name=tool-group]").Attr("value")
		if !exists {
			fmt.Println("No tool-group value found")
			return
		}

		// URL encode the toolGroup
		toolGroupEncoded := url.QueryEscape(toolGroup)
		urStr := fmt.Sprintf("/product-search/catalog/item/?tool-group=%s&abbreviation_method=partial_match&page_no=1", toolGroupEncoded)
		urlCollections = append(urlCollections, ninjacrawler.UrlCollection{
			Url:      ctx.App.GetFullUrl(urStr),
			MetaData: nil,
			Parent:   ctx.UrlCollection.Url,
		})
	})
	return urlCollections
}
