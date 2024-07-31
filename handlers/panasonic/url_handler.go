package panasonic

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

func UrlHandler(crawler *ninjacrawler.Crawler) {

	categorySelector := ninjacrawler.UrlSelector{
		Selector:     "ul li.cmp-list__item",
		FindSelector: "a.c-product__link",
		Attr:         "href",
		Handler: func(urlCollection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {

			if strings.Contains(fullUrl, "#") {
				return "", nil
			}
			return fullUrl, nil
		},
	}
	crawler.CrawlUrls([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.Categories,
			OriginCollection: crawler.GetBaseCollection(),
			Processor:        categorySelector,
		},
		{
			Entity:           constant.Products,
			OriginCollection: constant.Categories,
			Processor:        handleProduct,
		},
	})

}
func isValidHost(urlString string) bool {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Url parsing error:", err)
		return false
	}

	hostname := parsedUrl.Hostname()
	if hostname == "panasonic.jp" || hostname == "ec-plus.panasonic.jp" {
		return true
	}

	return false
}
func handleProduct(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	items := []ninjacrawler.UrlCollection{}

	ctx.Document.Find("a.normal").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok || len(href) == 0 || strings.Contains(href, "#") || strings.Contains(href, "javascript") {
			fmt.Println("Invalid href:", href)
			return
		}
		if !strings.HasPrefix(href, "/") {
			href = "/" + href
		}
		href = ctx.App.GetFullUrl(href)
		if !isValidHost(href) {
			return
		}

		items = append(items, ninjacrawler.UrlCollection{
			Url:    href,
			Parent: ctx.UrlCollection.Url,
		})
	})
	return items

}
