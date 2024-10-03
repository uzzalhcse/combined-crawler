package kojima

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type Category struct {
	Child []Category `json:"child"`
	Url   string     `json:"url"`
}

func UrlHandler(crawler *ninjacrawler.Crawler) {

	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categoryHandler,
		},
		{
			Entity:           constant.SubCategories,
			OriginCollection: constant.Categories,
			Processor:        subCategoryHandler,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.SubCategories,
			Processor:        productHandler,
		},
	})
}

func categoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	urls := []ninjacrawler.UrlCollection{}
	strContent := ctx.Document.Find("#MK2HEAD_CATE").Text()
	// Parse the JSON
	var categories []Category
	err := json.Unmarshal([]byte(strContent), &categories)
	if err != nil {
		fmt.Println("Error parsing")
	}

	// Extract and print child categories
	for _, category := range categories {
		for _, child := range category.Child {
			urls = append(urls, ninjacrawler.UrlCollection{
				Url:    ctx.App.GetFullUrl(child.Url),
				Parent: ctx.UrlCollection.Url,
			})
		}
	}
	return urls
}

func subCategoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	subCatUrls := []ninjacrawler.UrlCollection{}
	RecursiveSubCategoryCrawler(ctx, ctx.Document, &subCatUrls, ctx.UrlCollection.Url)
	return subCatUrls
}
func RecursiveSubCategoryCrawler(ctx ninjacrawler.CrawlerContext, doc *goquery.Document, subCatUrls *[]ninjacrawler.UrlCollection, urlStr string) {
	subCategoryList := doc.Find("ul#ChangToProdUrl > li")
	if subCategoryList.Length() > 1 {
		subCategoryList.Each(func(i int, s *goquery.Selection) {
			href, ok := s.Find("a").First().Attr("href")
			if ok {
				fullUrl := ctx.App.GetFullUrl(href)
				var document *goquery.Document

				navigationContext, err := ctx.App.Navigate(fullUrl)
				if err != nil {
					_ = ctx.App.MarkAsMaxErrorAttempt(ctx.UrlCollection.Url, constant.Categories, err.Error())
					ctx.App.Logger.Error("Error navigating to sub-category page : %v", err)
				}
				document = navigationContext.Document

				//fmt.Println("fullUrl:", fullUrl)
				RecursiveSubCategoryCrawler(ctx, document, subCatUrls, fullUrl) // Recursive call
			}
		})
	} else {
		liTags := ctx.Document.Find("li[pn='stock']")
		if liTags.Length() > 0 && urlStr != "" {
			*subCatUrls = append(*subCatUrls, ninjacrawler.UrlCollection{Url: urlStr, Parent: ctx.UrlCollection.Url})
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
