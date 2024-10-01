package ninjacrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func (app *Crawler) handleCrawlWorker(processorConfig ProcessorConfig, proxy Proxy, urlCollection UrlCollection) (*CrawlerContext, error) {
	//var page playwright.Page
	//var browser playwright.Browser
	var err error
	var doc *goquery.Document
	var apiResponse map[string]interface{}
	if *app.engine.IsDynamic {
		if *app.engine.Adapter == PlayWrightEngine {
			//browser, page, err = app.GetBrowserPage(app.pw, app.engine.BrowserType, proxy)
			//if err != nil {
			//	app.Logger.Fatal(err.Error())
			//}
			//defer browser.Close()
			//defer page.Close()
		} else {
			if app.rodBrowser == nil {
				fmt.Println("ROD BROWSER IS NIL")
			}
			if app.rodPage == nil {
				fmt.Println("ROD PAGE IS NIL")
			}

		}
	}

	if *app.engine.Adapter == PlayWrightEngine {
		page, pError := app.GetPage(app.browser)
		if pError != nil {
			app.Logger.Fatal(pError.Error())
		}
		app.page = page
		defer app.page.Close()
	} else {
		rodPage, err := app.GetRodPage(app.rodBrowser)
		if err != nil {
			app.Logger.Fatal(err.Error())
		}

		app.rodPage = rodPage
		defer app.rodPage.MustClose()
	}

	crawlableUrl := urlCollection.Url
	if urlCollection.ApiUrl != "" {
		crawlableUrl = urlCollection.ApiUrl
	}
	if urlCollection.CurrentPageUrl != "" {
		crawlableUrl = urlCollection.CurrentPageUrl
	}
	if proxy.Server != "" {
		app.Logger.Info("Crawling :%s: %s using Proxy %s", processorConfig.OriginCollection, crawlableUrl, proxy.Server)
	} else {
		app.Logger.Info("Crawling :%s: %s", processorConfig.OriginCollection, crawlableUrl)
	}
	if *app.engine.IsDynamic {
		if *app.engine.Adapter == PlayWrightEngine {
			doc, err = app.NavigateToURL(app.page, crawlableUrl)
		} else {
			doc, err = app.NavigateRodURL(app.rodPage, crawlableUrl)
		}
	} else {
		switch processorConfig.Processor.(type) {
		case ProductDetailApi:
			apiResponse, err = app.NavigateToApiURL(app.httpClient, crawlableUrl, proxy)
		default:
			doc, err = app.NavigateToStaticURL(app.httpClient, crawlableUrl, proxy)
		}
	}

	if err != nil {
		return nil, err
	}
	crawlerCtx := &CrawlerContext{
		App:           app,
		Document:      doc,
		UrlCollection: urlCollection,
		Page:          app.page,
		RodPage:       app.rodPage,
		ApiResponse:   apiResponse,
	}
	return crawlerCtx, nil
}
