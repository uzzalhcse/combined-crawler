package sandvik

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"time"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {

	crawler.Collection(constant.Categories).SetBrowserType("firefox").CrawlUrls(crawler.GetBaseCollection(), handleCategory)

}

func handleCategory(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	items, err := ctx.Page.Locator(".position-relative.search-push-wrapper.ng-star-inserted").All()
	if err != nil {
		ctx.App.Logger.Info("Error fetching items:", err)
		return urls
	}
	ctx.App.Logger.Info("Total Items: ", len(items))
	ctx.App.Logger.Html(ctx.Page, "hudai")

	for _, item := range items {
		time.Sleep(time.Second * 5)
		err := item.Click()
		if err != nil {
			ctx.App.Logger.Error("Failed to click on Product Card: %v", err)
		}

		time.Sleep(time.Second * 5)
		_, err = ctx.Page.GoBack()
		if err != nil {
			ctx.App.Logger.Error("Failed to goback: %v", err)
		}
	}
	return urls
}
