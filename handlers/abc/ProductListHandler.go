package main

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func ProductListHandler(urlCollection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	fmt.Println("Im from custom plugin ProductListHandler.", fullUrl)
	return fullUrl, nil
}
