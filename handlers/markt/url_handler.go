package markt

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"time"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	crawler.Collection(constant.Categories).CrawlUrls(crawler.GetBaseCollection(), ninjacrawler.UrlSelector{
		Selector:     ".l-category-button-list__in",
		SingleResult: false,
		FindSelector: "a.c-category-button",
		Attr:         "href",
	})
	crawler.Collection(constant.Products).CrawlUrls(constant.Categories, handleProducts)
}

func handleProducts(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	productLinkSelector := ".p-product-detail__review a.c-text-link.u-color-text--link.c-text-link--underline"
	clickAndWaitButton(ctx.App, ".u-hidden-sp li button", ctx.Page)
	const maxRetries = 1
	const retryInterval = 10 * time.Second
	items, err := ctx.Page.Locator("ul.p-card-list-no-scroll li.p-product-card.p-product-card--large").All()
	if err != nil {
		ctx.App.Logger.Info("Error fetching items:", err)
		return urls
	}

	for i, item := range items {
		err = retryWithSleep(maxRetries, retryInterval, func(attempt int) error {
			if attempt > 1 {
				ctx.App.Logger.Info("Attempt #%d, waiting %v before retrying item index %d of %d & URL %s", attempt, retryInterval, i, len(items), ctx.UrlCollection.Url)
			}
			err := item.Click()
			if err != nil {
				ctx.App.Logger.Error("Failed to click on Product Card: %v", err)
				return err
			}

			// Wait for the modal to open and the link to be available
			_, err = ctx.Page.WaitForSelector(productLinkSelector)
			if err != nil {
				closeModal(ctx)
				ctx.App.Logger.Warn("Failed To Open Modal, item index %d of %d %s %v", i, len(items), ctx.UrlCollection.Url, err.Error())
				return err
			}

			doc, err := ctx.App.GetPageDom(ctx.Page)
			if err != nil {
				ctx.App.Logger.Error("Error getting page DOM:", err)
				return err
			}

			productLink, exist := doc.Find(productLinkSelector).First().Attr("href")
			fullUrl := ctx.App.GetFullUrl(productLink)
			if !exist {
				ctx.App.Logger.Error("Failed to find product link")
				return fmt.Errorf("product link not found")
			}
			urls = append(urls, ninjacrawler.UrlCollection{Url: fullUrl, Parent: ctx.UrlCollection.Url})

			return nil
		})

		if err != nil {
			ctx.App.Logger.Error("Error processing item: %v", err)
			closeModal(ctx)
			continue
		}

		closeModal(ctx)

		// Add a delay after every 50 items
		if (i+1)%50 == 0 {
			ctx.App.Logger.Info("Sleeping for 5 seconds...")
			time.Sleep(5 * time.Second)
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
		ctx.App.Logger.Error("WaitFor Close Modal %v", err)
	}
	closeModal := ctx.Page.Locator("#__next > div.l-background__wrap > div.l-background__in > div > button")
	if closeModal != nil {
		err = closeModal.Click(playwright.LocatorClickOptions{Timeout: playwright.Float(10000)})
		if err != nil {
			ctx.App.Logger.Warn("Failed To Close Modal %s %v", ctx.UrlCollection.Url, err.Error())
		}

	} else {
		ctx.App.Logger.Error("Modal close button not found.")
	}
	_, err = ctx.Page.WaitForSelector("l-background__wrap", playwright.PageWaitForSelectorOptions{
		State: playwright.WaitForSelectorStateDetached,
	})
	if err != nil {
		ctx.App.Logger.Error("WaitForSelectorStateDetached %v", err)
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
