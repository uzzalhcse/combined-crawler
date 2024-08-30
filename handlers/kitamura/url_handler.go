package kitamura

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	categorySelector := ninjacrawler.UrlSelector{
		Selector:     ".category-item",
		FindSelector: "a",
		Attr:         "href",
	}
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categorySelector,
		},
	})

}

func subCategoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}

	// Process the DOM elements if they exist
	ctx.Document.Find("ul.category_list li,ul li").Each(func(_ int, s *goquery.Selection) {
		s.Find("div.category_order h4 a").Each(func(_ int, a *goquery.Selection) {
			href, exists := a.Attr("href")
			if exists {
				brand1 := a.Text()
				href = ctx.App.GetFullUrl(href)
				if !isValidHost(href) {
					return
				}
				items = append(items, GetProductUrls(ctx, href, brand1)...)
			}
		})
	})

	return items
}
func GetProductUrls(ctx ninjacrawler.CrawlerContext, subCategory string, brand1 string) []ninjacrawler.UrlCollection {
	var productUrls []ninjacrawler.UrlCollection

	doc, err := ctx.App.NavigateToURL(ctx.Page, subCategory)
	if err != nil {
		_ = ctx.App.MarkAsError(ctx.UrlCollection.Url, constant.Categories, err.Error())
		return productUrls
	}
	// Check if "ul.category_list li" exists in the DOM
	if doc.Find(".category_wrap").Length() == 0 {
		return []ninjacrawler.UrlCollection{{
			Url:    subCategory,
			Parent: ctx.UrlCollection.Url,
			MetaData: map[string]interface{}{
				"brand1": brand1,
			},
		}}
	}
	doc.Find("ul.category_list li,ul li").Each(func(i int, li *goquery.Selection) {
		li.Find("div.category_order h4 a").Each(func(j int, a *goquery.Selection) {
			brand2 := a.Text()
			href, exists := a.Attr("href")
			if exists {
				href = ctx.App.GetFullUrl(href)
				if !isValidHost(href) {
					return
				}
				productUrls = append(productUrls, ninjacrawler.UrlCollection{
					Url:    href,
					Parent: subCategory,
					MetaData: map[string]interface{}{
						"brand1": brand1,
						"brand2": brand2,
					},
				})
			}
		})
	})

	return productUrls
}
