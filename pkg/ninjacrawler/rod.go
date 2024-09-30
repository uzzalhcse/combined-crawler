package ninjacrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"time"
)

// GetRodBrowser initializes and runs Rod browser.
// It returns a Rod browser instance if successful, otherwise returns an error.
func (app *Crawler) GetRodBrowser(proxy Proxy) (*rod.Browser, error) {
	// Setup the browser launcher with proxy if provided
	l := launcher.New().Headless(!app.isLocalEnv).Devtools(app.isLocalEnv)

	if len(app.engine.ProxyServers) > 0 && proxy.Server != "" {
		l = l.Set(flags.ProxyServer, proxy.Server)
	}

	url := l.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	// Optionally handle proxy authentication
	if proxy.Username != "" && proxy.Password != "" {
		go browser.MustHandleAuth(proxy.Username, proxy.Password)()
	}

	return browser, nil
}

// GetRodPage creates a new page using the Rod framework.
// It returns the page instance, or an error if the operation fails.
func (app *Crawler) GetRodPage(browser *rod.Browser) (*rod.Page, error) {
	page := browser.MustPage()

	err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: app.userAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("error setting user agent: %s", err.Error())
	}

	return page, nil
}

// NavigateRodURL navigates to a specified URL using the Rod page.
// It waits until the page is fully loaded, handles cookie consent, and returns the page DOM.
func (app *Crawler) NavigateRodURL(page *rod.Page, url string) (*goquery.Document, error) {
	timeout := time.Duration(app.engine.Timeout) * time.Second

	// Go to the URL with a timeout
	err := page.Timeout(timeout).Navigate(url)
	if err != nil {
		//app.Logger.Html(page.MustHTML(), url, err.Error())
		return nil, err
	}

	// Optionally wait for a specific selector
	if app.engine.WaitForSelector != nil {
		elem, _ := page.Timeout(time.Duration(app.engine.Timeout) * time.Second).Element(*app.engine.WaitForSelector)
		if elem == nil {
			fmt.Printf("For URL: %s, element not found!\n")
			page.MustWaitStable()
		}
	} else {
		page.MustWaitLoad()
	}

	// Handle cookie consent
	err = app.HandleRodCookieConsent(page)
	if err != nil {
		app.Logger.Html(page.MustHTML(), url, err.Error())
		return nil, err
	}

	// Get the page DOM
	document, err := app.GetRodPageDom(page)

	// Optionally send HTML to BigQuery or store it
	if app.engine.SendHtmlToBigquery != nil && *app.engine.SendHtmlToBigquery {
		sendErr := app.SendHtmlToBigquery(document, url)
		if sendErr != nil {
			app.Logger.Fatal("SendHtmlToBigquery Error: %s", sendErr.Error())
		}
	}
	if *app.engine.StoreHtml {
		if err := app.SaveHtml(document, url); err != nil {
			app.Logger.Error(err.Error())
		}
	}

	return document, nil
}

// HandleRodCookieConsent handles cookie consent dialogs on the page.
func (app *Crawler) HandleRodCookieConsent(page *rod.Page) error {
	action := app.engine.CookieConsent
	if action == nil {
		return nil
	}

	// Fill input fields if specified
	for _, field := range action.Fields {
		el := page.MustElement(fmt.Sprintf("input[name='%s']", field.Key))
		el.MustInput(field.Val)
	}

	// Click the consent button
	if action.ButtonText != "" {
		button := page.MustElementR("button", action.ButtonText)
		button.MustClick()

		page.MustWaitLoad()
	}

	return nil
}
