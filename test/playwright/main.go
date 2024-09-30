package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
)

func main() {
	// Start Playwright in non-headless mode for debugging
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	// Launch browser in headful mode for visual debugging
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // Set to false to see the browser in action
		Devtools: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	// Set the User-Agent
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"

	// Create a new browser context with the custom User-Agent
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String(userAgent),
		//JavaScriptEnabled: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not create new browser context: %v", err)
	}
	defer context.Close()

	// Define the maximum resource loading time (30 seconds)
	//maxResourceLoadTime := 60 * time.Second
	//
	//// Intercept and monitor network requests
	//context.Route("**/*", func(route playwright.Route) {
	//	startTime := time.Now()
	//	request := route.Request()
	//	// Allow request to continue and asynchronously check the response time
	//	go func() {
	//		// Wait for the resource to load or timeout
	//		for {
	//			if time.Since(startTime) > maxResourceLoadTime {
	//				log.Printf("Request took too long: %s - Aborting it.", request.URL())
	//				//route.Abort()
	//				//return
	//			}
	//
	//			// Check if the request has finished
	//			if request.Timing().ResponseEnd > 0 {
	//				break
	//			}
	//
	//			// Sleep a little before checking again
	//			time.Sleep(100 * time.Millisecond)
	//		}
	//		// Continue the request if it did not timeout
	//		route.Continue()
	//	}()
	//})

	// Open a new page with the modified context
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create new page: %v", err)
	}

	// Navigate to the target URL
	url := "https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NI-U701-K"
	_, err = page.Goto(url, playwright.PageGotoOptions{
		Timeout: playwright.Float(60000), // Increase timeout to 60 seconds
	})
	if err != nil {
		log.Fatalf("could not navigate to page: %v", err)
	}
	// Add delay for visual inspection (only for debugging)
	//time.Sleep(10 * time.Second)
	// Wait for product detail selector (increase timeout)
	//_, err = page.WaitForSelector("h1.pd_c-headingLv1-01", playwright.PageWaitForSelectorOptions{
	//	//Timeout: playwright.Float(60000), // Increase timeout for waiting on this selector
	//})
	//if err != nil {
	//	log.Fatalf("could not find product detail section: %v", err)
	//}

	// Extract some information from the page (e.g., product name)
	productName, err := page.TextContent("h1.pd_c-headingLv1-01")
	if err != nil {
		log.Fatalf("could not extract product name: %v", err)
	}
	fmt.Printf("Product Name: %s\n", productName)

	// Example of extracting other data (e.g., price)
	price, err := page.TextContent("div.pd_c-price")
	if err != nil {
		log.Fatalf("could not extract product price: %v", err)
	}
	fmt.Printf("Price: %s\n", price)

	// Close the browser and cleanup
	page.Close()
	browser.Close()
	pw.Stop()
}
