package yamaya

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func productNameHandler(ctx ninjacrawler.CrawlerContext) string {
	productName := strings.Trim(ctx.Document.Find("div.card-header.bg-white").First().Text(), " \n")

	return productName
}

func getUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}

func getProductCategory(ctx ninjacrawler.CrawlerContext) string {
	var categoryTexts []string

	categoryDiv := ctx.Document.Find("div.card")
	categoryDiv.Each(func(i int, s *goquery.Selection) {
		cardBody := s.Find("div.card-body").Last()
		tableTag := cardBody.Find("table.table.table-sm.table-borderless")
		aTag := tableTag.Find("a").First()
		text := aTag.Text()
		categoryTexts = append(categoryTexts, strings.TrimSpace(text))
	})

	return strings.Join(categoryTexts, " > ")
}

func getProductDescription(ctx ninjacrawler.CrawlerContext) string {
	description := ctx.Document.Find(".card-body .card-text:not(.text-right)").Text()
	return description
}
func getSellingPrice(ctx ninjacrawler.CrawlerContext) string {
	priceText := ctx.Document.Find("button.btnAddToCart span").Text()

	reg := regexp.MustCompile(`[^0-9]`)
	priceText = reg.ReplaceAllString(priceText, "")
	return priceText

}
func getProductAttribute(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}

	ctx.Document.Find("div.card-body table tr").Each(func(i int, s *goquery.Selection) {
		th := strings.Trim(s.Find("th").Text(), " \n")
		td := strings.Trim(s.Find("td").Text(), " \n")

		if len(th) > 0 && len(td) > 0 {
			if th != "分類" {
				attribute := ninjacrawler.AttributeItem{
					Key:   th,
					Value: td,
				}
				attributes = append(attributes, attribute)
			}
		}
	})

	return attributes
}
