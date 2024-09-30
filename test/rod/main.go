package main

import (
	"fmt"
	"github.com/go-rod/rod/lib/proto"
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	// Create a launcher with headless mode (or disable headless mode by setting the flag to false)
	url, _ := launcher.New().Headless(false).Launch()

	// Create a new browser instance with the launched URL
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()
	UserAgent := &proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	}
	// Create a new page with a custom User-Agent
	page := browser.MustPage().MustSetUserAgent(UserAgent)

	// Navigate to the URL
	productUrl := "https://ec-plus.panasonic.jp/store/ap/storeaez/a2A/ProductDetail?HB=NI-U701-K"
	if err := page.Navigate(productUrl); err != nil {
		log.Fatalf("could not navigate to the URL: %v", err)
	}
	cssConditionalElement := "h1.pd_c-headingLv1-01"
	elem, _ := page.Timeout(time.Duration(30) * time.Second).Element(cssConditionalElement)
	if elem == nil {
		fmt.Printf("For URL: %s, element not found!\n")
		page.MustWaitStable()
	}
	// Extract product name (Assuming it's inside an h1 tag with class product-title)
	productName, _ := page.Element("h1.pd_c-headingLv1-01")
	productNameText, _ := productName.Text()
	fmt.Printf("Product Name: %s\n", productNameText)

	// Extract product price (Assuming it's inside a span tag with class price)
	price := page.MustElement("div.pd_c-price").MustText()
	fmt.Printf("Price: %s\n", price)

	// Close the browser
	page.MustClose()
}
