package sandvik

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/playwright-community/playwright-go"
	"strings"
	"time"
)

const (
	baseUrl   = "https://www.sandvik.coromant.com"
	sleepTime = 5
)

func productHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	productTableSelector := ".mb-1.mb-md-2.mt-2.flex-grow-1"
	loadMoreButtonSelector := ".row.mt-2.ng-star-inserted .ng-star-inserted"
	ctx.Page.WaitForSelector(productTableSelector)

	err := ctx.Page.Click(loadMoreButtonSelector)
	if err != nil {
		ctx.App.Logger.Error(err.Error())
		return nil
	}
	time.Sleep(time.Second * sleepTime) //slow down to avoid hitting the rate limit

	jsScroll := `
			   var intervalID = setInterval(function () {
				  var scrollingElement = (document.scrollingElement || document.body);
				  var footer = document.querySelector('.footer-container.cor-color-regular-white');
				  var footerHeight = footer ? footer.offsetHeight : 0;
				  scrollingElement.scrollTop = scrollingElement.scrollHeight - scrollingElement.clientHeight - footerHeight;
			   }, 200);
			`
	_, err = ctx.Page.Evaluate(jsScroll)
	if err != nil {
		ctx.App.Logger.Error("Error setting up scroll interval: %v", err)
		return nil
	}
	scrollSelector := ".row.mt-2.ng-star-inserted"

	for {
		count, err := ctx.Page.Locator(scrollSelector).Count()
		if err != nil {
			ctx.App.Logger.Error("Error checking selector existence: %v", err)
			break
		}
		if count == 0 {
			ctx.App.Logger.Info("Selector found, breaking scroll loop")
			break
		}
		ctx.Page.WaitForSelector(".row.mt-2.ng-star-inserted .col-12 .ng-star-inserted", playwright.PageWaitForSelectorOptions{
			State: playwright.WaitForSelectorStateDetached,
		})
	}

	productList, err := ctx.Page.Locator(".flex-grow-1.flex-md-grow-0 a").All()
	if err != nil {
		ctx.App.Logger.Error(err.Error())
	}
	for _, product := range productList {
		url, err := product.GetAttribute("href")
		if err != nil {
			ctx.App.Logger.Error(err.Error())
			continue
		}
		textToFind := "https://"
		if !strings.Contains(url, textToFind) {
			url = baseUrl + url
			urls = append(urls, ninjacrawler.UrlCollection{Url: url, Parent: ctx.UrlCollection.Url})
		}
	}
	return urls
}

func categoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	time.Sleep(sleepTime * time.Second) //slow down to avoid hitting the rate limit
	items, err := ctx.Page.Locator(".position-relative.search-push-wrapper.ng-star-inserted").All()
	if err != nil {
		ctx.App.Logger.Info("Error fetching items:", err)
		return urls
	}

	for _, item := range items {
		ctx.Page.WaitForSelector(".position-relative.search-push-wrapper.ng-star-inserted")
		err := item.Click()
		if err != nil {
			ctx.App.Logger.Error("Failed to click on Product Card: %v", err)
		}
		ctx.Page.WaitForSelector(".mb-1.mb-md-2.mt-2.flex-grow-1")

		urls = append(urls, ninjacrawler.UrlCollection{
			Url: ctx.Page.URL(),
		})
		_, err = ctx.Page.GoBack()
		time.Sleep(sleepTime * time.Second) //slow down to avoid hitting the rate limit
		if err != nil {
			ctx.App.Logger.Error("Failed to go back: %v", err)
		}
	}
	return urls
}
