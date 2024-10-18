package sony

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func ProductNameHandler(ctx ninjacrawler.CrawlerContext) string {
	return strings.Trim(ctx.Document.Find(".details .intro h2").First().Text(), " \n")
}

func GetUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}
func GetProductCodes(ctx ninjacrawler.CrawlerContext) []string {
	var codes []string
	ctx.Document.Find(".ProductSummaryModels__ModelCode").Each(func(i int, s *goquery.Selection) {
		// Skip the first two items
		codes = append(codes, s.Text())
	})
	return codes
}

func GetProductDescription(ctx ninjacrawler.CrawlerContext) string {

	description := ctx.Document.Find(".details .intro .text p").Text()
	description = strings.ReplaceAll(description, "\n\n", "\n")

	return description
}
func GetProductAttribute(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}

	GetSpecAttributeService(ctx, &attributes)
	GetMeritAttributeService(ctx.App, ctx.Document, &attributes)
	GetCatalogAttributeService(ctx.App, ctx.Document, &attributes)

	return attributes
}

func GetSpecAttributeService(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	var specUrl string

	// Define a helper function to find the spec URL
	findSpecUrl := func(selector, text string) string {
		var url string
		ctx.Document.Find(selector).EachWithBreak(func(i int, s *goquery.Selection) bool {
			if strings.Contains(s.Text(), text) {
				url, _ = s.Attr("href")
				return false // Break the loop as soon as we find a match
			}
			return true
		})
		return url
	}

	// Try different logic in sequence to find the spec URL
	specUrl = findSpecUrl("ul.CategoryNav__PdpNavList li.CategoryNav__PdpNavItem a.CategoryNav__PdpNavItemLink", "主な仕様")
	if specUrl == "" {
		specUrl = findSpecUrl("a.s5-buttonV3", "すべての仕様を見る")
	}
	if specUrl == "" && ctx.Document.Find(".s5-specTable").Length() > 0 {
		// If the spec table exists in the DOM, we assume this is the page we need
		fmt.Println("Found spec table")
		specUrl = ctx.UrlCollection.Url
	}

	// Handle the case where no spec URL was found
	if specUrl == "" {
		ctx.App.Logger.Error("Failed to get spec url %s", ctx.UrlCollection.Url)
		_ = ctx.App.MarkAsMaxErrorAttempt(ctx.UrlCollection.Url, constant.Products, "Failed to get spec url")
		return
	}
	//if !strings.HasPrefix(specUrl, "/") {
	//	specUrl = "/" + specUrl
	//}
	// If the spec URL is not absolute, make it absolute
	if !strings.HasPrefix(specUrl, "http://") && !strings.HasPrefix(specUrl, "https://") {
		specUrl = ctx.UrlCollection.Url + specUrl
	}

	ctx.App.Logger.Warn("Spec Url: %s", specUrl)

	// Navigate to the spec URL
	_, err := ctx.App.NavigateToStaticURL(ctx.App.GetHttpClient(), specUrl, ctx.App.CurrentProxy)
	if err != nil {
		ctx.App.Logger.Error("Failed to navigate to spec url %s: %v", specUrl, err)
		return
	}
}

func GetMeritAttributeService(app *ninjacrawler.Crawler, document *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	key := strings.Trim(document.Find(".merit.clearfix h3").First().Text(), " \n")
	values := strings.Trim(document.Find(".merit.clearfix ul").First().Text(), " \n")

	if len(values) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: values,
		}
		*attributes = append(*attributes, attribute)
	}
}

func GetCatalogAttributeService(app *ninjacrawler.Crawler, document *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	document.Find("#detail ul li").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a")
		key := strings.Trim(a.Text(), " \n")
		img := s.Find("img")
		alt, exist := img.Attr("alt")
		if exist {
			key = alt
		}

		value, exists := a.Attr("href")

		if exists {
			fullUrl := app.GetFullUrl(value)

			attribute := ninjacrawler.AttributeItem{
				Key:   key,
				Value: fullUrl,
			}
			*attributes = append(*attributes, attribute)
		}
	})
}
