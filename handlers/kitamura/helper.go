package kitamura

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

func getJanService(ctx ninjacrawler.CrawlerContext) string {
	var janCode string
	ctx.Document.Find("dl#product_detail_standard span.product_detail_item,.product_detail dl span").Each(func(
		i int, s *goquery.Selection) {
		dt := s.Find("dt").Text()
		if dt == "JANコード" {
			janCode = s.Find("dd").Text()
		}
	})
	return janCode
}

func GetUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}
func getProductNameService(ctx ninjacrawler.CrawlerContext) string {
	productName := ctx.Document.Find("div#product,div#product_detail").Find("h2").Text()

	return productName
}

func getCategoryService(ctx ninjacrawler.CrawlerContext) string {
	var categoryTexts []string
	ctx.Document.Find("span.gt").Remove()
	categoryDiv := ctx.Document.Find(".topicpath.pc_only > *:not(span.gt)")
	categoryDiv.Each(func(i int, a *goquery.Selection) {
		text := a.Text()
		categoryTexts = append(categoryTexts, strings.TrimSpace(text))
	})

	category := strings.Join(categoryTexts, " > ")

	return category
}

func getDescriptionService(ctx ninjacrawler.CrawlerContext) string {
	// Replace <br> tags with newline characters
	ctx.Document.Find("p#product_detail_exp br,.product_detail p br").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml("\n")
	})

	// Extract the text content
	descriptionText := ctx.Document.Find("p#product_detail_exp,.product_detail p").Text()

	// Trim leading and trailing whitespace
	return strings.TrimSpace(descriptionText)
}
func getListPriceService(ctx ninjacrawler.CrawlerContext) string {
	var listPrice string
	ctx.Document.Find("dl#product_detail_standard span.product_detail_item,.product_detail dl span").Each(func(
		i int, s *goquery.Selection) {
		dt := s.Find("dt").Text()
		if dt == "希望小売価格" {
			listPrice = strings.TrimSuffix(s.Find("dd").Text(), "円")
		}
	})

	return listPrice
}

func getProductAttribute(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}
	if ctx.UrlCollection.MetaData["brand1"] != nil {
		attributes = append(attributes, ninjacrawler.AttributeItem{
			Key:   "brand1",
			Value: ctx.UrlCollection.MetaData["brand1"].(string),
		})
	}
	if ctx.UrlCollection.MetaData["brand2"] != nil {
		attributes = append(attributes, ninjacrawler.AttributeItem{
			Key:   "brand2",
			Value: ctx.UrlCollection.MetaData["brand2"].(string),
		})
	}
	productName := getProductNameService(ctx)
	getTagsAttributeService(productName, &attributes)
	getReleaseDateAttributeService(ctx.Document, &attributes)
	getKeyAttributeService(ctx.Document, &attributes)
	return attributes
}
func getTagsAttributeService(productName string, attributes *[]ninjacrawler.AttributeItem) {
	tag := ""
	if strings.Contains(productName, "特定保健用食品") {
		tag = "特定保健用食品"
	} else if strings.Contains(productName, "機能性表示食品") {
		tag = "機能性表示食品"
	}
	if len(tag) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "tags",
			Value: tag,
		}
		*attributes = append(*attributes, attribute)
	}
}
func getReleaseDateAttributeService(doc *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	text := doc.Find("p#product_detail_exp").Find("span.new").Text()

	if len(text) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "発売日",
			Value: text,
		}
		*attributes = append(*attributes, attribute)
	}
}

func getKeyAttributeService(doc *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	// product_details keys-values
	var capacity, expiry string
	doc.Find("dl#product_detail_standard span.product_detail_item,.product_detail dl span").Each(func(
		i int, s *goquery.Selection) {
		dt := s.Find("dt").Text()
		if strings.TrimSpace(dt) == "容量" {
			capacity = s.Find("dd").Text()
		} else if strings.TrimSpace(dt) == "賞味期間" {
			expiry = s.Find("dd").Text()
		}
	})

	if len(capacity) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "容量",
			Value: capacity,
		}
		*attributes = append(*attributes, attribute)
	}

	if len(expiry) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "賞味期間",
			Value: expiry,
		}
		*attributes = append(*attributes, attribute)
	}

	// product_block2 keys-values
	doc.Find("div#product_block2").Find("dl").Find("span.product_block2_group").Each(func(
		i int, s *goquery.Selection) {
		s.Find("span.product_block2_item").Each(func(j int, span *goquery.Selection) {
			dt := span.Find("dt").First().Text()
			dd := span.Find("dd")

			if len(dt) > 0 {
				dlText := dd.Find("dl").Text()
				if len(dlText) > 0 {
					ddText := ""
					dd.Find("dl").Find("span").Each(func(k int, sp *goquery.Selection) {
						dt_ := sp.Find("dt").Text()
						dd_ := sp.Find("dd").Text()
						if dt_ != "" && dd_ != "" {
							ddText += dt_ + " " + dd_ + " \n"
						}
					})
					attribute := ninjacrawler.AttributeItem{
						Key:   strings.TrimSpace(dt),
						Value: strings.TrimSpace(ddText),
					}
					*attributes = append(*attributes, attribute)
				} else {
					attribute := ninjacrawler.AttributeItem{
						Key:   strings.TrimSpace(dt),
						Value: strings.TrimSpace(dd.Text()),
					}
					*attributes = append(*attributes, attribute)
				}
			}
		})

	})
}

func isValidHost(urlString string) bool {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Url parsing error:", err)
		return false
	}

	hostname := parsedUrl.Hostname()
	if hostname == "suntory.co.jp" || hostname == "products.suntory.co.jp" {
		return true
	}

	return false
}
