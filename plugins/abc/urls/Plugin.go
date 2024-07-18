package main

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func HandleCategoryUrl(ctx ninjacrawler.CrawlerContext) []ninjacrawler.UrlCollection {
	var urls []ninjacrawler.UrlCollection
	ctx.Document.Find("ul.Header_Navigation_List_Item_Sub_Group_Inner").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			attrValue, ok := s.Attr("href")
			if !ok {
				ctx.App.Logger.Error("Attribute not found. %v", "href")
			} else {
				fullUrl := ctx.App.GetFullUrl(attrValue)
				urls = append(urls, ninjacrawler.UrlCollection{Url: fullUrl, Parent: ctx.UrlCollection.Parent})
				fmt.Println("Im from custom plugin HandleCategory.", fullUrl)
			}
		})
	})
	return urls
}
func ProductListHandler(urlCollection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	fmt.Println("Im from custom plugin ProductListHandler.", fullUrl)
	return fullUrl, nil
}
