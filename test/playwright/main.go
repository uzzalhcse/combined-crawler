package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"sync"
)

func crawlURL(context playwright.BrowserContext, url string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()

	// Open a new tab (page) within the same browser context
	page, err := context.NewPage()
	if err != nil {
		log.Printf("could not create new page for %s: %v", url, err)
		<-semaphore
		return
	}
	defer page.Close()

	// Navigate to the target URL
	_, err = page.Goto(url, playwright.PageGotoOptions{
		Timeout: playwright.Float(60000), // 60-second timeout
	})
	if err != nil {
		log.Printf("====could not navigate to page %s: %v", url, err)
		<-semaphore
		return
	}

	// Extract product name
	productName, err := page.TextContent("h1.pd_c-headingLv1-01")
	if err != nil {
		log.Printf("could not extract product name from %s: %v", url, err)
		<-semaphore
		return
	}
	fmt.Printf("Product Name from %s: %s\n", url, productName)

	// Extract price
	price, err := page.TextContent("div.pd_c-price")
	if err != nil {
		log.Printf("could not extract product price from %s: %v", url, err)
		<-semaphore
		return
	}
	fmt.Printf("Price from %s: %s\n", url, price)

	// Release the semaphore slot to allow other goroutines to proceed
	<-semaphore
}

func main() {
	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	// Launch a single browser instance
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // Set to false to see the browser in action
		Devtools: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	// Set User-Agent for the browser context
	userAgent := "PostmanRuntime/7.37.3"
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String(userAgent),
	})
	if err != nil {
		log.Fatalf("could not create new browser context: %v", err)
	}
	defer context.Close()

	// Define the URLs to crawl
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

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) // Limit concurrency to 3 tabs at a time

	for _, url := range urls {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire a semaphore slot
		go crawlURL(context, url, &wg, semaphore)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Crawling completed!")
}
