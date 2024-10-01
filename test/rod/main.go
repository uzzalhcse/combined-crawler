package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func main() {
	urls := []string{
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NI-U701-K",
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NI-A66-K",
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NI-FS70A-K",
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NA-LX129DR-W",
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NA-FA11K3-N",
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NA-FA12V3-W",
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NA-FA8H3-W",
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=HH-CL1492A",
		"https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=SQ-LD440-W",
	}

	// Create a launcher with headless mode (or disable headless mode by setting the flag to false)
	url, _ := launcher.New().Headless(true).Launch()

	// Create a new browser instance with the launched URL
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// Wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a buffered channel to limit concurrency (3 concurrent goroutines)
	semaphore := make(chan struct{}, 3)

	// Crawl each URL concurrently with a limit of 3 at a time
	for _, productUrl := range urls {
		wg.Add(1) // Increment the wait group counter

		// Acquire a slot in the semaphore before launching a goroutine
		semaphore <- struct{}{}

		go func(url string) {
			defer wg.Done()                // Decrement the counter when the goroutine completes
			defer func() { <-semaphore }() // Release the slot in the semaphore

			// Crawl product details
			crawlProduct(browser, url)
		}(productUrl)
	}

	// Wait for all crawls to finish
	wg.Wait()
}

func crawlProduct(browser *rod.Browser, productUrl string) {
	// Create a new page with a custom User-Agent
	UserAgent := &proto.NetworkSetUserAgentOverride{
		UserAgent: "PostmanRuntime/7.37.3",
	}
	page := browser.MustPage().MustSetUserAgent(UserAgent)

	// Navigate to the URL
	if err := page.Navigate(productUrl); err != nil {
		log.Fatalf("could not navigate to the URL: %v", err)
		return
	}

	// Wait for the product name element to appear
	cssConditionalElement := "h1.pd_c-headingLv1-01"
	elem, _ := page.Timeout(30 * time.Second).Element(cssConditionalElement)
	if elem == nil {
		fmt.Printf("For URL: %s, element not found!\n", productUrl)
		page.MustWaitStable()
		return
	}

	// Extract product name
	productName, _ := elem.Text()
	fmt.Printf("Product Name: %s\n", productName)

	// Extract product price
	price := page.MustElement("div.pd_c-price").MustText()
	fmt.Printf("Price for %s: %s\n", productName, price)

	// Close the page
	page.MustClose()
}
