package panasonic

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func getImagesService(ctx ninjacrawler.CrawlerContext) []string {
	baseUrl := ctx.App.BaseUrl
	images := []string{}

	url := strings.ReplaceAll(ctx.UrlCollection.Url, baseUrl, "")
	urlStarting := strings.Split(url, "/")[1]

	mainSection := ctx.Document.Find("section#maincontents").First()
	if mainSection == nil {
		return images
	}
	productDetailsSections := getProductDetailsSections(ctx.Document)
	sections := []*goquery.Selection{}
	firstFigure := mainSection.Find("figure").First()

	sections = append(sections, firstFigure)
	sections = append(sections, productDetailsSections)

	uniqueImages := make(map[string]bool)
	for _, section := range sections {
		if section == nil {
			continue
		}
		section.Find("img").Each(func(i int, s *goquery.Selection) {
			dataSrc, dataSrcOk := s.Attr("data-src")
			src, _ := s.Attr("src")
			imgUrl := ""

			if dataSrcOk {
				imgUrl = dataSrc
			} else {
				imgUrl = src
			}
			if !strings.HasPrefix(imgUrl, "/"+urlStarting) {
				return
			}
			url := baseUrl + imgUrl

			if !uniqueImages[url] {
				uniqueImages[url] = true
				images = append(images, url)
			}
		})
	}

	return images
}
func getProductCodesService(ctx ninjacrawler.CrawlerContext) []string {
	productCodes := []string{}
	if !strings.Contains(ctx.UrlCollection.Url, "products/") {
		return productCodes
	}

	productUrlParts := strings.Split(ctx.UrlCollection.Url, "/")
	productCode := productUrlParts[len(productUrlParts)-1]
	productCode = strings.Split(productCode, ".")[0]

	productCodes = append(productCodes, productCode)

	return productCodes
}
func ProductNameHandler(ctx ninjacrawler.CrawlerContext) string {

	productMainSection := getProductMainSection(ctx.Document)
	productName := ""
	if productMainSection == nil {
		return productName
	}

	productNameSection := productMainSection.Find("h1").First()
	productName = strings.Trim(productNameSection.Text(), " \t\n")

	return productName
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

	description := ctx.Document.Find("p.product_handler-header-lead").Text()
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
	item := strings.Trim(document.Find(".product_handler-header-category").First().Text(), " \n")

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
	key := document.Find(".product_handler-feature-inner .section-title-ja.typesquare_option").First().Text()
	value := document.Find(".product_handler-feature-lead").Text()
	document.Find(".product_handler-feature-section").Each(func(i int, s *goquery.Selection) {
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
	key := document.Find(".product_handler-material-title").First().Text()
	values := make([]string, 0)
	document.Find(".modal-content .product_handler-material ul.product_handler-material-list li").Each(func(i int, s *goquery.Selection) {
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

func getProductDetailsSections(pageData *goquery.Document) *goquery.Selection {
	productDetailsSectionHeaders := pageData.Find("div.bgLightGray")

	var productDetailsSections *goquery.Selection = nil
	productDetailsSectionHeaders.Each(func(i int, s *goquery.Selection) {
		h2Text := s.Find("h2").First().Text()
		text := strings.Trim(h2Text, " \n")
		if text == "新商品のおすすめポイント" || text == "ご購入をお考えのお客様へ" {
			return
		}

		selection := s
		for {
			selection = selection.Next()
			if selection.Text() == "" || selection.HasClass("bgLightGray") {
				break
			}
			if productDetailsSections == nil {
				productDetailsSections = selection
			} else {
				productDetailsSections = productDetailsSections.AddSelection(selection)
			}
		}
	})

	return productDetailsSections
}
func getProductMainSection(pageData *goquery.Document) *goquery.Selection {
	mainSection := pageData.Find("section#maincontents").First()

	var productMainSection *goquery.Selection
	mainSection.Find("div.pagesection").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("bgLightGray") {
			return
		}
		if productMainSection == nil {
			productMainSection = s
		} else {
			productMainSection = productMainSection.AddSelection(s)
		}
	})

	return productMainSection
}
