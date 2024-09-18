package topvalu

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"time"
)

const (
	baseUrl   = "https://www.topvalu.net"
	sleepTime = 1
)

func CategoryHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	categoryDiv := ctx.Document.Find("div.pulldown__list--head").First()
	categoryDiv.Find("div.pulldown__inner--parent").Each(func(i int, s *goquery.Selection) {
		linkElement := s.Find("a.pulldown__ttl").First()
		href, ok := linkElement.Attr("href")
		if !ok {
			ctx.App.Logger.Error("category url not found.")
			return
		}
		urls = append(urls, ninjacrawler.UrlCollection{Url: baseUrl + href, Parent: ctx.UrlCollection.Url})
	})
	return urls
}

func ProductHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	productUrlSelector := ".product.product--category .product__item .product__item__title a"
	loadMoreButton := ".product__button"
	click, _ := ctx.Page.Locator(loadMoreButton).Count()
	for click != 0 {
		err := ctx.Page.Locator(loadMoreButton).Click()
		if err != nil {
			ctx.App.Logger.Error("load more button click error.", err)
		}
		time.Sleep(sleepTime * time.Second)
		click, _ = ctx.Page.Locator(loadMoreButton).Count()
	}
	count, err := ctx.Page.Locator(productUrlSelector).Count()
	if err != nil {
		return nil
	}
	ctx.App.Logger.Info("product count:", count)

	productCards, err := ctx.Page.Locator(productUrlSelector).All()
	if err != nil {
		ctx.App.Logger.Error(err.Error())
	}
	for _, productCard := range productCards {
		href, _ := productCard.GetAttribute("href")
		ctx.App.Logger.Info("product url:", baseUrl+href)
		urls = append(urls, ninjacrawler.UrlCollection{Url: baseUrl + href, Parent: ctx.UrlCollection.Url})
	}
	return urls
}
