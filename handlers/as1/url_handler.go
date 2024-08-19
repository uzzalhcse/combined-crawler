package as1

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
		},
	})
}
func categoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	categoryDiv := ctx.Document.Find("#menuList > div.submenu > div.accordion > dl")

	// Get the total number of items in the selection
	totalCats := categoryDiv.Length()

	// Iterate over all items except the last three
	categoryDiv.Slice(0, totalCats-3).Each(func(i int, cat *goquery.Selection) {
		cat.Find("dd > div > ul > li").Each(func(j int, li *goquery.Selection) {
			li.Find("ul > li").Each(func(j int, l *goquery.Selection) {
				l.Find("ul > li").Each(func(j int, lMain *goquery.Selection) {
					href, ok := lMain.Find("a").Attr("href")
					if ok {
						urls = append(urls, ninjacrawler.UrlCollection{
							Url:    ctx.App.GetFullUrl(href),
							Parent: ctx.UrlCollection.Url,
						})
						ctx.App.Logger.Info("Category URL %s", href)
					} else {
						ctx.App.Logger.Error("Category URL not found")
					}
				})
			})
		})
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
