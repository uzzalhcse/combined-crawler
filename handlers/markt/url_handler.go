package markt

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/playwright-community/playwright-go"
	"time"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor: ninjacrawler.UrlSelector{
				Selector:     ".l-category-button-list__in",
				SingleResult: false,
				FindSelector: "a.c-category-button",
				Attr:         "href",
			},
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        handleProducts,
		},
	})
}

func handleProducts(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	productLinkSelector := ".p-product-detail__review a.c-text-link.u-color-text--link.c-text-link--underline"
	clickAndWaitButton(ctx.App, ".u-hidden-sp li button", ctx.Page)

	items, err := ctx.Page.Locator("ul.p-card-list-no-scroll li.p-product-card.p-product-card--large").All()
	if err != nil {
		ctx.App.Logger.Warn("Error fetching items:", err)
		return urls
	}

	for i, item := range items {
		links, _ := item.Locator("a.p-product-card__wrap").Count()
		modals, _ := item.Locator("div.p-product-card__wrap").Count()
		if links > 0 {
			attribute, err := item.Locator("a.p-product-card__wrap").GetAttribute("href")
			if err != nil {
				ctx.App.Logger.Warn("Failed to Get Attribute", err)
				continue
			}

			fullUrl := ctx.App.GetFullUrl(attribute)
			urls = append(urls, ninjacrawler.UrlCollection{Url: fullUrl, Parent: ctx.UrlCollection.Url})
		} else if modals > 0 {
			if err := item.Click(); err != nil {
				ctx.App.Logger.Warn("Failed to click on Product Card: %v", err)
				continue
			}

			if _, err := ctx.Page.WaitForSelector(productLinkSelector); err != nil {
				ctx.App.Logger.Warn("Failed to open modal, item index %d of %d %s %v", i, len(items), ctx.UrlCollection.Url, err)
				closeModal(ctx)
				continue
			}

			doc, err := ctx.App.GetPageDom(ctx.Page)
			if err != nil {
				ctx.App.Logger.Warn("Error getting page DOM:", err)
				closeModal(ctx)
				continue
			}

			productLink, exist := doc.Find(productLinkSelector).First().Attr("href")
			if !exist {
				ctx.App.Logger.Warn("Failed to find product link")
				closeModal(ctx)
				continue
			}

			fullUrl := ctx.App.GetFullUrl(productLink)
			urls = append(urls, ninjacrawler.UrlCollection{Url: fullUrl, Parent: ctx.UrlCollection.Url})
			closeModal(ctx)
		}

	}

	return urls
}

// retryWithSleep retries the given function fn up to maxRetries times with the specified sleep interval between retries.
func retryWithSleep(maxRetries int, sleepInterval time.Duration, fn func(attempt int) error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn(i + 1)
		if err == nil {
			return nil
		}
		time.Sleep(sleepInterval)
	}
	return err
}
func closeModal(ctx ninjacrawler.CrawlerContext) {
	// Close the modal

	//ctx.Page.Locator("l-modal-content l-modal-content--frame-low")
	_, err := ctx.Page.WaitForSelector("#__next > div.l-background__wrap > div.l-background__in > div > button")
	if err != nil {
		ctx.App.Logger.Warn("WaitFor Close Modal %v", err)
	}
	closeModal := ctx.Page.Locator("#__next > div.l-background__wrap > div.l-background__in > div > button")
	if closeModal != nil {
		err = closeModal.Click(playwright.LocatorClickOptions{Timeout: playwright.Float(10000)})
		if err != nil {
			ctx.App.Logger.Warn("Failed To Close Modal %s %v", ctx.UrlCollection.Url, err.Error())
		}

	} else {
		ctx.App.Logger.Warn("Modal close button not found.")
	}
	_, err = ctx.Page.WaitForSelector("l-background__wrap", playwright.PageWaitForSelectorOptions{
		State: playwright.WaitForSelectorStateDetached,
	})
	if err != nil {
		ctx.App.Logger.Warn("WaitForSelectorStateDetached %v", err)
	}
}
func clickAndWaitButton(crawler *ninjacrawler.Crawler, selector string, page playwright.Page) {
	for {
		button := page.Locator(selector)
		err := button.Click()
		page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(1000)})
		if err != nil {
			crawler.Logger.Info("No more button available")
			break
		}
	}
}
