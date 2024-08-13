package kojima

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"time"
)

type Category struct {
	Child []Category `json:"child"`
	Url   string     `json:"url"`
}

func UrlHandler(crawler *ninjacrawler.Crawler) {

	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
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
	fmt.Println("total subCatUrls", len(subCatUrls))
	return subCatUrls
}
func RecursiveSubCategoryCrawler(ctx ninjacrawler.CrawlerContext, doc *goquery.Document, subCatUrls *[]ninjacrawler.UrlCollection, urlStr string) {
	subCategoryList := doc.Find("ul#ChangToProdUrl > li")
	if subCategoryList.Length() > 1 {
		subCategoryList.Each(func(i int, s *goquery.Selection) {
			href, ok := s.Find("a").First().Attr("href")
			if ok {
				fullUrl := ctx.App.GetFullUrl(href)
				httpClient := ctx.App.GetHttpClient()

				var document *goquery.Document
				var err error

				maxAttempts := 3
				for attempt := 1; attempt <= maxAttempts; attempt++ {
					document, err = ctx.App.NavigateToStaticURL(httpClient, fullUrl, ctx.App.CurrentProxy)
					if err == nil {
						break // Successful navigation, exit retry loop
					}

					ctx.App.Logger.Warn("Attempt %d: Error navigating to sub-category page: %v", attempt, err)

					if attempt == maxAttempts {
						_ = ctx.App.MarkAsError(ctx.UrlCollection.Url, constant.Categories, err.Error(), 1)
						ctx.App.Logger.Error("Error navigating to sub-category page after %d attempts: %v", maxAttempts, err)
						return
					}
					time.Sleep(5 * time.Second)
				}

				fmt.Println("fullUrl:", fullUrl)
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
