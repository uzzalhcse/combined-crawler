package kojima

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	categorySelector := ninjacrawler.UrlSelector{
		Selector:     "ul.MK2PFRDH000_01 li li",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
	}

	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categorySelector,
			Engine:           ninjacrawler.Engine{IsDynamic: true},
		},
		{
			Entity:           constant.SubCategories,
			OriginCollection: constant.Categories,
			Processor:        subCategoryHandler,
			Engine:           ninjacrawler.Engine{IsDynamic: true},
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.SubCategories,
			Processor:        productHandler,
			Engine:           ninjacrawler.Engine{IsDynamic: true},
		},
	})
}

func subCategoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	subCatUrls := []ninjacrawler.UrlCollection{}
	RecursiveSubCategoryCrawler(ctx, ctx.Document, &subCatUrls)
	return subCatUrls
}
func RecursiveSubCategoryCrawler(ctx ninjacrawler.CrawlerContext, doc *goquery.Document, subCatUrls *[]ninjacrawler.UrlCollection) {

	subCategoryList := doc.Find("ul#ChangToProdUrl > li")

	if subCategoryList.Length() > 1 {
		subCategoryList.Each(func(i int, s *goquery.Selection) {
			href, ok := s.Find("a").First().Attr("href")
			if ok {
				fullUrl := ctx.App.GetFullUrl(href)
				document, err := ctx.App.NavigateToURL(ctx.Page, fullUrl)
				if err != nil {
					ctx.App.Logger.Error("Error navigating to sub-category page:", err)
					return
				}
				RecursiveSubCategoryCrawler(ctx, document, subCatUrls) // Recursive call
			}
		})
	} else {
		liTags := ctx.Document.Find("li[pn='stock']")
		if liTags.Length() > 0 {
			*subCatUrls = append(*subCatUrls, ninjacrawler.UrlCollection{Url: ctx.Page.URL(), Parent: ctx.UrlCollection.Url})
			return
		} else {
			return
		}
	}
}
func productHandler(ctx ninjacrawler.CrawlerContext, fn func([]ninjacrawler.UrlCollection, string)) error {

	items := []ninjacrawler.UrlCollection{}
	ctx.Document.Find(".name a.mk2TagClick").Each(func(_ int, s *goquery.Selection) {
		productLink, exists := s.Attr("href")
		if exists {
			items = append(items, ninjacrawler.UrlCollection{
				Url:    ctx.App.GetFullUrl(productLink),
				Parent: ctx.UrlCollection.Url,
			})
		}
	})
	nextPage := ctx.Document.Find("a.next")
	nextPageUrl, _ := nextPage.Attr("href")
	if nextPage.Length() == 0 {
		fn(items, "")
		return nil
	} else {
		fn(items, ctx.App.GetFullUrl(nextPageUrl))
	}
	return nil
}
