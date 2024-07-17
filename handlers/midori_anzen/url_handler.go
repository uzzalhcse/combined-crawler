package midori_anzen

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categoryHandler,
			//Preference:       ninjacrawler.Preference{DoNotMarkAsComplete: true},
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        productHandler,
			Engine: ninjacrawler.Engine{
				IsDynamic: true,
			},
			Preference: ninjacrawler.Preference{DoNotMarkAsComplete: true},
		},
	})

}
func categoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	ctx.Document.Find("select.form-control.form-control-search option").Each(func(i int, s *goquery.Selection) {
		attrValue, ok := s.Attr("value")
		if ok {
			if attrValue != "" {
				fmt.Println("attrValue", attrValue)
				pageUrl := fmt.Sprintf("/shop/goods/search.aspx?tree=%s&keyword=&search=x", attrValue)
				fullUrl := ctx.App.GetFullUrl(pageUrl)
				urls = append(urls, ninjacrawler.UrlCollection{Url: fullUrl, Parent: ctx.UrlCollection.Url})
			}
		}
	})
	return urls
}
func productHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	currentPage := "1"

	for {
		items, err := ctx.Page.Locator(".pagination.pagination-sm li").All()
		if err != nil {
			fmt.Printf("Failed to locate pagination items: %v\n", err)
			return urls
		}
		if len(items) < 2 {
			fmt.Println("Pagination items are insufficient")
			return urls
		}

		lastPage := items[len(items)-2]
		lastPageNumber, err := lastPage.TextContent()
		if err != nil {
			fmt.Printf("Failed to get last page number: %v\n", err)
			return urls
		}

		fmt.Printf("c %s, l %s\n", currentPage, lastPageNumber)
		if currentPage == lastPageNumber {
			fmt.Println("No more next page")
			break
		}

		for _, item := range items {
			urls = append(urls, getUrls(ctx)...)

			currentPage, err = item.Locator("a").TextContent()
			if err != nil {
				fmt.Printf("Failed to get current page number: %v\n", err)
				return urls
			}

			fmt.Printf("currentPage %s, lastPageNumber %s\n", currentPage, lastPageNumber)
			if currentPage == lastPageNumber {
				fmt.Println("No more next page...")
				break
			}

			err = item.Locator("a").Click()
			if err != nil {
				fmt.Printf("Failed to click on next page: %v\n", err)
				return urls
			}

			ctx.Page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
				State: playwright.LoadStateNetworkidle,
			})
		}
	}
	return urls
}

func getUrls(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	items, err := ctx.Page.Locator(".event-goods .iconcard.event-price-img").All()
	if err != nil {
		ctx.App.Logger.Warn("Error fetching items:", err)
		return urls
	}

	for _, item := range items {

		attribute, err := item.Locator("a").GetAttribute("href")
		if err != nil {
			ctx.App.Logger.Warn("Failed to Get Attribute", err)
			continue
		}

		fullUrl := ctx.App.GetFullUrl(attribute)
		urls = append(urls, ninjacrawler.UrlCollection{Url: fullUrl, Parent: ctx.UrlCollection.Url})
	}

	return urls
}
