package osg

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func getUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	if ctx.UrlCollection.CurrentPageUrl != "" {
		return ctx.UrlCollection.CurrentPageUrl
	}
	return ctx.UrlCollection.Url
}

func getSellingPrice(ctx ninjacrawler.CrawlerContext) string {
	priceText := ctx.Document.Find("button.btnAddToCart span").Text()

	reg := regexp.MustCompile(`[^0-9]`)
	priceText = reg.ReplaceAllString(priceText, "")
	return priceText

}
func getProductAttribute(ctx ninjacrawler.CrawlerContext, selection *goquery.Selection) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}

	selection.Find("dd button.btn").Remove()

	selection.Find("details div dl.clearfix").Each(func(i int, s *goquery.Selection) {
		s.Find("dt").Each(func(j int, dt *goquery.Selection) {
			key := dt.Find("label").Text()
			val := strings.Trim(dt.Next().Text(), " \n") // Get the next sibling which is the corresponding `dd`

			if len(key) > 0 && len(val) > 0 {
				if j > 0 {
					attribute := ninjacrawler.AttributeItem{
						Key:   key,
						Value: val,
					}
					attributes = append(attributes, attribute)
				}
			}
		})
	})

	// Extract key-value pairs from p elements while excluding certain classes
	selection.Find("div.productimage > p").Not(".thumb, .wide, .favorite").Each(func(i int, p *goquery.Selection) {
		a := p.Find("a")
		if a.Length() > 0 {
			key := a.Text()
			val := a.AttrOr("href", "")
			if len(key) > 0 && len(val) > 0 {
				attribute := ninjacrawler.AttributeItem{
					Key:   key,
					Value: ctx.App.GetFullUrl(val),
				}
				attributes = append(attributes, attribute)
			}
		}
	})

	// Handle DXF and STEP separately to get file names
	selection.Find("div.productimage p.download").Each(func(i int, p *goquery.Selection) {
		p.Find("label a").Each(func(j int, a *goquery.Selection) {
			key := a.Parent().AttrOr("aria-label", "")
			val := a.AttrOr("href", "")
			if len(key) > 0 && len(val) > 0 {
				attribute := ninjacrawler.AttributeItem{
					Key:   key,
					Value: ctx.App.GetFullUrl(val),
				}
				attributes = append(attributes, attribute)
			}
		})
	})

	return attributes
}
