package main

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func CategoryHandler(urlCollection ninjacrawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	fmt.Println("Im from custom plugin CategoryHandler.", fullUrl)
	return fullUrl, nil
}
