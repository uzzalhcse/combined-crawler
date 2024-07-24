package main

import (
	"combined-crawler/pkg/generic_crawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func Categories(ctx generic_crawler.CrawlerContext, next func([]generic_crawler.UrlCollection, string)) error {
	var urls []generic_crawler.UrlCollection
	ctx.Document.Find("ul.nav li.nav-item").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			attrValue, ok := s.Attr("href")
			if !ok {
				fmt.Println("Attribute not found. %v", "href")
			} else {
				fullUrl := ctx.App.GetFullUrl(attrValue)
				urls = append(urls, generic_crawler.UrlCollection{Url: fullUrl, Parent: ctx.UrlCollection.Parent})
				fmt.Println("Im from custom plugin HandleCategory.", fullUrl)
			}
		})
	})
	next(urls, "")
	return nil
}
func Products(urlCollection generic_crawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	fmt.Println("Im from custom plugin ProductListHandler.", fullUrl)
	return fullUrl, nil
}
