package rakuten

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func getJanService(ctx ninjacrawler.CrawlerContext) string {
	janCode := ""
	if ctx.State.Get("jsonData.gtin13") != nil {
		janCode = ctx.State.Get("jsonData.gtin13").(string)
	}
	return janCode
}
func GetUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}

func getImagesService(ctx ninjacrawler.CrawlerContext) []string {
	var images []string
	if imgInterfaces, ok := ctx.State.Get("jsonData.image").([]interface{}); ok {
		for _, img := range imgInterfaces {
			if str, ok := img.(string); ok {
				images = append(images, str)
			}
		}
	}
	return images
}
func getProductCodesService(ctx ninjacrawler.CrawlerContext) []string {
	productCodes := []string{}

	if ctx.State.Get("jsonData.sku") != nil {
		productCodes = append(productCodes, ctx.State.Get("jsonData.sku").(string))
	}

	return productCodes
}
func getBrandService(ctx ninjacrawler.CrawlerContext) string {
	brand := ""

	if ctx.State.Get("jsonData.brand.name") != nil {
		brand = ctx.State.Get("jsonData.brand.name").(string)
	}

	return brand
}
func getProductNameService(ctx ninjacrawler.CrawlerContext) string {
	productName := ctx.Document.Find("h1.pd_c-headingLv1-01").First().Text()
	if len(productName) == 0 {
		productName = ctx.Document.Find("h1.pd_categorybox_product__heading").First().Text()
	}

	return productName
}
func GetProductCategory(ctx ninjacrawler.CrawlerContext) string {
	category := ""
	categoryList := []string{}

	ulItem := ctx.Document.Find("div.pd_c-breadcrumb > ul").First()
	ulItem.Find("li").Each(func(i int, s *goquery.Selection) {
		categoryText := strings.Trim(s.Text(), " \n")
		categoryList = append(categoryList, categoryText)
	})

	category = strings.Join(categoryList, " > ")

	return category
}

func GetProductDescription(ctx ninjacrawler.CrawlerContext) string {
	description := ""

	descriptionSection := ctx.Document.Find("p.pd_b-detail-02_text-01").First()
	descriptionSection.Find("style").Remove()
	descriptionSection.Find("a").Remove()

	description = ctx.App.HtmlToText(descriptionSection)
	description = strings.Split(description, "◇開梱・設置オプションをお申し込みのお客様へ")[0]
	description = strings.Split(description, "◆リサイクルオプションをお申し込みのお客様へ")[0]
	description = strings.Trim(description, " \n")

	return description
}
func getSellingPriceService(ctx ninjacrawler.CrawlerContext) string {
	sellingPrice := ""

	priceDivText := ctx.Document.Find("div.pd_c-price").First().Text()
	sellingPrice = ctx.App.ToNumericsString(priceDivText)

	return sellingPrice
}

func GetProductAttribute(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	var attributes []ninjacrawler.AttributeItem
	//specificationData := ctx.State.Get("specificationData").(map[string]string)
	//sellingPrice := getSellingPriceService(ctx)
	//
	//getURLRelationProductDetailsAttributeService(ctx.Document, &attributes)
	//getURLRelationProductSpecAttributeService(ctx.Document, &attributes)
	//getSellingPriceTaxAttributeService(sellingPrice, &attributes)
	//getSpecialKeyAttributeService(ctx.Document, specificationData, &attributes)

	return attributes
}

func getURLRelationProductDetailsAttributeService(pageData *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	productDetailsURL := ""

	productDetailsURLSection := pageData.Find("a.pd_m-linkList-01_anchor--features").First()
	if productDetailsURLSection.Length() == 0 {
		productDetailsURLSection = pageData.Find("p.pd_b-detail-02_text-01").First()
		if productDetailsURLSection.Length() != 0 {
			productDetailsURLSection = productDetailsURLSection.Find("a").First()
		}
	}

	productDetailsURL, ok := productDetailsURLSection.Attr("href")
	if ok {
		attribute := ninjacrawler.AttributeItem{
			Key:   "url_relation_product_details",
			Value: productDetailsURL,
		}
		*attributes = append(*attributes, attribute)
	}
}

func getURLRelationProductSpecAttributeService(pageData *goquery.Document, attributes *[]ninjacrawler.AttributeItem) {
	specificationSection := pageData.Find("a.pd_m-linkList-01_anchor--spec").First()
	specificationURL, ok := specificationSection.Attr("href")
	if ok {
		attribute := ninjacrawler.AttributeItem{
			Key:   "url_relation_product_spec",
			Value: specificationURL,
		}
		*attributes = append(*attributes, attribute)
	}
}

func getSellingPriceTaxAttributeService(sellingPrice string, attributes *[]ninjacrawler.AttributeItem) {
	if len(sellingPrice) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "selling_price_tax",
			Value: "1",
		}
		*attributes = append(*attributes, attribute)
	}
}

func getSpecialKeyAttributeService(pageData *goquery.Document, specificationData map[string]string, attributes *[]ninjacrawler.AttributeItem) {
	specialKeySection := pageData.Find("p.pd_c-released-01").First()
	if specialKeySection != nil {
		contents := specialKeySection.Text()
		contents = strings.Split(contents, "/")[0]
		contents = strings.Trim(contents, " \n")
		if len(contents) == 0 {
			return
		}

		specialKey := "発売日"
		contentsSlice := strings.Split(contents, "：")
		if len(contentsSlice) < 2 {
			return
		}
		key := strings.Trim(contentsSlice[0], " \n")
		value := strings.Trim(contentsSlice[1], " \n")
		if key == specialKey && len(value) > 0 {
			attribute := ninjacrawler.AttributeItem{
				Key:   specialKey,
				Value: value,
			}
			*attributes = append(*attributes, attribute)
		}
	}

	for key, value := range specificationData {
		attribute := ninjacrawler.AttributeItem{
			Key:   key,
			Value: value,
		}
		*attributes = append(*attributes, attribute)
	}
}
