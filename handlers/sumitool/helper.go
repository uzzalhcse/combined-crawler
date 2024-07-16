package sumitool

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func ProductNameHandler(ctx ninjacrawler.CrawlerContext) string {
	return strings.Trim(ctx.Document.Find("h1.product-header-name").First().Text(), " \n")
}

func GetUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}
func GetProductCategory(ctx ninjacrawler.CrawlerContext) string {
	categoryItems := make([]string, 0)
	// Find all <li> elements within .breadcrumb-nav ul, skipping the first two and the last item
	ctx.Document.Find(".breadclum-nav ul li:nth-child(n+3):not(:last-child)").Each(func(i int, s *goquery.Selection) {
		txt := strings.TrimSpace(s.Text())
		categoryItems = append(categoryItems, txt)
	})
	return strings.Join(categoryItems, " > ")
}

func GetProductDescription(ctx ninjacrawler.CrawlerContext) string {

	description := ctx.Document.Find("p.product-header-lead").Text()
	description = strings.ReplaceAll(description, "\n\n", "\n")

	return description
}
func GetProductAttribute(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}

	GetCatchCopyAttributeService(ctx.App, ctx.Document, &attributes)
	GetCatalogPDFAttributeService(ctx.App, ctx.Document, &attributes)
	GetFeatureAttributeService(ctx.App, ctx.Document, &attributes)
	GetMaterialAttributeService(ctx.App, ctx.Document, &attributes)

	return attributes
}

func GetCatchCopyAttributeService(app *ninjacrawler.Crawler, document *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	item := strings.Trim(document.Find(".product-header-category").First().Text(), " \n")

	if len(item) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "catch_copy",
			Value: item,
		}
		*attributes = append(*attributes, attribute)
	}
}
func GetCatalogPDFAttributeService(app *ninjacrawler.Crawler, document *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	ancor := document.Find("a.button.track_event.typesquare_option")
	key := ancor.Text()
	value, exist := ancor.Attr("href")
	if exist {
		fullUrl := app.GetFullUrl(value)
		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: fullUrl,
		}
		*attributes = append(*attributes, attribute)
	}
}

func GetFeatureAttributeService(app *ninjacrawler.Crawler, document *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	key := document.Find(".product-feature-inner .section-title-ja.typesquare_option").First().Text()
	value := document.Find(".product-feature-lead").Text()
	document.Find(".product-feature-section").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		// Replace multiple spaces/newlines with a single space
		re := regexp.MustCompile(`\s+`)
		cleanedText := re.ReplaceAllString(text, " ")
		value += "\n" + cleanedText + " "
	})

	if len(value) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: value,
		}
		*attributes = append(*attributes, attribute)
	}
}

func GetMaterialAttributeService(app *ninjacrawler.Crawler, document *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	key := document.Find(".product-material-title").First().Text()
	values := make([]string, 0)
	document.Find(".modal-content .product-material ul.product-material-list li").Each(func(i int, s *goquery.Selection) {
		s.Find("div").Remove()
		val := s.Text()
		values = append(values, strings.TrimSpace(val))
	})
	val := strings.Join(values, ",")
	if len(val) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: val,
		}
		*attributes = append(*attributes, attribute)
	}
}
