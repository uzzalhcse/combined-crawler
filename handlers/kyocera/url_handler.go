package kyocera

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	productSelector := ninjacrawler.UrlSelector{
		Selector:     "ul.heightLineParent.clearfix li",
		FindSelector: "dl dt a",
		Attr:         "href",
	}
	productOtherSelector := ninjacrawler.UrlSelector{
		Selector:     "ul.product-list li.product-item",
		FindSelector: "a",
		Attr:         "href",
	}
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        handleCategory,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        productSelector,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Other,
			Processor:        productOtherSelector,
		},
	})
}
func handleCategory(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	categoryUrls := []ninjacrawler.UrlCollection{}
	otherUrls := []ninjacrawler.UrlCollection{}
	productUrls := []ninjacrawler.UrlCollection{}
	ctx.Document.Find("div.index.clearfix ul.clearfix li").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		fullUrl := ctx.App.GetFullUrl(href)
		if strings.Contains(fullUrl, "/prdct/tool/category/product/") {
			categoryUrls = append(categoryUrls, ninjacrawler.UrlCollection{
				Url:    fullUrl,
				Parent: ctx.UrlCollection.Url,
			})
		} else if strings.Contains(fullUrl, "/prdct/tool/sgs/") {
			otherUrls = append(otherUrls, ninjacrawler.UrlCollection{
				Url:    fullUrl,
				Parent: ctx.UrlCollection.Url,
			})
		} else if strings.Contains(fullUrl, "/prdct/tool/product/") {
			productUrls = append(productUrls, ninjacrawler.UrlCollection{
				Url:    fullUrl,
				Parent: ctx.UrlCollection.Url,
			})
		}
	})
	ctx.App.InsertUrlCollections(constant.Other, otherUrls, ctx.UrlCollection.Url)
	ctx.App.InsertUrlCollections(constant.Products, productUrls, ctx.UrlCollection.Url)
	return categoryUrls
}
