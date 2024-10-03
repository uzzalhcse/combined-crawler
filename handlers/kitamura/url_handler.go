package kitamura

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {
	categorySelector := ninjacrawler.UrlSelector{
		Selector:     ".category-item",
		FindSelector: "a",
		Attr:         "href",
		Handler: func(urlCollection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
			fullUrl = fullUrl + "?limit=100&page=1"
			return fullUrl, nil
		},
	}
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categorySelector,
			Engine: ninjacrawler.Engine{
				WaitForSelector: ninjacrawler.String(".category-item"),
			},
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        productHandler,
			Engine: ninjacrawler.Engine{
				WaitForSelector: ninjacrawler.String("div#product-list-area>.product"),
			},
		},
	})

}

func productHandler(ctx ninjacrawler.CrawlerContext, next func([]ninjacrawler.UrlCollection, string)) error {
	productUrls := []ninjacrawler.UrlCollection{}
	productCountInfo := strings.TrimSpace(ctx.Document.Find("div.result-count").Text())

	re := regexp.MustCompile(`^\d+`)
	productCountStr := re.FindString(productCountInfo)
	if productCountStr == "" {
		return fmt.Errorf("could not find product count")
	}

	productCountInt, err := strconv.Atoi(productCountStr)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	totalPageNumber := int(math.Ceil(float64(productCountInt) / 100))

	currentPage := 1
	crawlableUrl := ctx.UrlCollection.Url
	if ctx.UrlCollection.CurrentPageUrl != "" {
		crawlableUrl = ctx.UrlCollection.CurrentPageUrl
	}

	parsedURL, err := url.Parse(crawlableUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return err
	}

	queryParams := parsedURL.Query()

	currentPageStr := queryParams.Get("page")
	if currentPageStr != "" {
		currentPage, err = strconv.Atoi(currentPageStr)
		if err != nil {
			fmt.Println("Error converting page to int:", err)
			return err
		}
	}
	productDiv := ctx.Document.Find("div#product-list-area").First()
	productDiv.Find("a.product-link").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			href = ctx.App.GetFullUrl(href)
			productUrls = append(productUrls, ninjacrawler.UrlCollection{
				Url:      href,
				MetaData: nil,
				Parent:   ctx.UrlCollection.Url,
			})
		} else {
			ctx.App.Logger.Warn("Product URL not found for %s", ctx.UrlCollection.Url)
		}

	})
	if currentPage < totalPageNumber {
		queryParams.Set("page", strconv.Itoa(currentPage+1))
		parsedURL.RawQuery = queryParams.Encode()
		nextPageUrl := parsedURL.String()
		fmt.Println("nextPageUrl:", nextPageUrl)
		next(productUrls, nextPageUrl)
	} else {
		next(productUrls, "")
	}

	return nil
}
