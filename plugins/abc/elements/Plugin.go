package main

import (
	"combined-crawler/pkg/generic_crawler"
	"strings"
)

func Url(ctx generic_crawler.CrawlerContext) interface{} {
	return ctx.UrlCollection.Url
}
func Name(ctx generic_crawler.CrawlerContext) interface{} {
	productName := ctx.Document.Find(".product-name h2").Text()
	productName = strings.Trim(productName, " \n")

	return productName
}

func Category(ctx generic_crawler.CrawlerContext) interface{} {
	category := ctx.Document.Find("p.ProductDetail_Section_Headline_Sub").First().Text()
	category = strings.Trim(category, " \n")

	return category
}

func Description(ctx generic_crawler.CrawlerContext) interface{} {
	description := ctx.Document.Find(".description p").Text()
	return description
}

func Attributes(ctx generic_crawler.CrawlerContext) interface{} {
	attributes := []generic_crawler.AttributeItem{}
	getExampleAttributeService(ctx, &attributes)
	return attributes
}

func getExampleAttributeService(ctx generic_crawler.CrawlerContext, attributes *[]generic_crawler.AttributeItem) {
	item := strings.Trim(ctx.Document.Find(".example p").First().Text(), " \n")
	if len(item) > 0 {
		attribute := generic_crawler.AttributeItem{
			Key:   "example",
			Value: item,
		}
		*attributes = append(*attributes, attribute)
	}
}
