package topvalu

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"path/filepath"
	"regexp"
	"strings"
)

func ProductDetailsHandler(crawler *ninjacrawler.Crawler) {
	ProductDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan:              JanHandler,
		PageTitle:        pageTitleHandler,
		Url:              urlHandler,
		Images:           imageHandler,
		ProductCodes:     productCodeHandler,
		Maker:            "",
		Brand:            brandHandler,
		ProductName:      productNameHandler,
		Category:         productCategoryHandler,
		Description:      productDescriptionHandler,
		Reviews:          reviewsHandler,
		ItemTypes:        func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		ItemSizes:        func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		ItemWeights:      func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		SingleItemSize:   SingleItemSizeHandler,
		SingleItemWeight: SingleItemWeightHandler,
		NumOfItems:       "",
		ListPrice:        "",
		SellingPrice:     sellingPriceHandler,
		Attributes:       productAttributesHandler,
	}

	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        ProductDetailSelector,
			Preference:       ninjacrawler.Preference{ValidationRules: []string{"PageTitle"}},
			Engine: ninjacrawler.Engine{
				BlockResources: false,
			},
		},
	})

}

func pageTitleHandler(ctx ninjacrawler.CrawlerContext) string {
	pageTitle := ctx.Document.Find("title").Text()
	return pageTitle
}

func urlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}

func JanHandler(ctx ninjacrawler.CrawlerContext) string {
	return getSpecificationData(ctx, "JAN")
}

func SingleItemSizeHandler(ctx ninjacrawler.CrawlerContext) string {
	return getSpecificationData(ctx, "サイズ")
}

func SingleItemWeightHandler(ctx ninjacrawler.CrawlerContext) string {
	return getSpecificationData(ctx, "重さ")
}

func productCodeHandler(ctx ninjacrawler.CrawlerContext) []string {
	return []string{""}
}

func brandHandler(ctx ninjacrawler.CrawlerContext) string {
	var brand string
	brandSection := ctx.Document.Find("a.item-detail__label--brand").First()
	brandImage := brandSection.Find("img").First()
	alt, ok := brandImage.Attr("alt")
	if ok {
		brand = alt
	}
	brand = strings.ReplaceAll(brand, "ロゴ画像", "")
	return brand
}

func productNameHandler(ctx ninjacrawler.CrawlerContext) string {
	productName := ctx.Document.Find("h1.item-detail__name").First().Text()
	return productName
}

func productCategoryHandler(ctx ninjacrawler.CrawlerContext) string {
	var categories []string
	ctx.Document.Find("li.breadcrumb__item").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		categories = append(categories, s.Text())
	})
	category := strings.Join(categories, " > ")
	return category
}

func productDescriptionHandler(ctx ninjacrawler.CrawlerContext) string {
	var description string
	descriptionDiv := ctx.Document.Find("div.item-detail__informations").First()
	descriptionDiv.Find("p").Each(func(i int, s *goquery.Selection) {
		description += s.Text()
		description += "\n"
	})
	description = strings.Trim(description, "\n")
	return description
}

func getSpecificationData(ctx ninjacrawler.CrawlerContext, ind string) string {
	data := make(map[string]string)
	var keys []string
	var values []string

	ctx.Document.Find("dt.item-detail__specs__field").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		text = strings.ReplaceAll(text, "：", "")
		keys = append(keys, text)
	})
	ctx.Document.Find("dd.item-detail__specs__value").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		values = append(values, text)
	})

	ctx.Document.Find("div.infoBox").Each(func(i int, s *goquery.Selection) {
		s.Find("dt").Each(func(j int, ss *goquery.Selection) {
			keys = append(keys, ss.Text())
		})
		s.Find("dd").Each(func(j int, ss *goquery.Selection) {
			values = append(values, ss.Text())
		})
	})

	allergySection := ctx.Document.Find("div#link_allergy").First()
	allergySectionHeader := allergySection.Next().Text()
	allergySectionHeader = strings.Trim(allergySectionHeader, "\n")
	allergySectionHeader = strings.Trim(allergySectionHeader, "\t")
	allergySectionHeader = strings.Trim(allergySectionHeader, "\n")

	var allergyItems []string
	ctx.Document.Find("li.items-allergy__list__item").Each(func(i int, s *goquery.Selection) {
		caption := s.Find("figcaption").First().Text()
		img := s.Find("img").First()
		src, ok := img.Attr("src")
		if ok {
			fileName := filepath.Base(src)
			splitted := strings.Split(fileName, "_")
			if splitted[1] == "on" {
				allergyItems = append(allergyItems, caption)
			}
		}
	})
	allergyItemsValue := strings.Join(allergyItems, ",")
	keys = append(keys, allergySectionHeader)
	values = append(values, allergyItemsValue)

	for index, key := range keys {
		data[key] = values[index]
	}
	if _, ok := data[ind]; ok {
		return data[ind]
	} else {
		return ""
	}
}

func sellingPriceHandler(ctx ninjacrawler.CrawlerContext) string {
	priceSection := ctx.Document.Find("dd.item-detail__specs__value--price").First()
	priceSection = priceSection.Find("strong.item-detail__specs__strong").First()
	priceText := priceSection.Text()

	//price := scripts.GetPriceWithSymbol(priceText)
	reg := regexp.MustCompile(`[^0-9]`)
	str := reg.ReplaceAllString(priceText, "")
	return str
}

func imageHandler(ctx ninjacrawler.CrawlerContext) []string {
	var imageList []string

	imageUl := ctx.Document.Find("ul.item-detail__thumbs").First()
	imageUl.Find("li.item-detail__thumb").Each(func(i int, s *goquery.Selection) {
		img := s.Find("img").First()
		src, ok := img.Attr("src")
		if ok {
			imageList = append(imageList, src)
		}
	})

	return imageList
}

func productAttributesHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	var attributes []ninjacrawler.AttributeItem
	getSellingPriceTaxAttributeService(ctx, &attributes)
	getSpecialKeyAttributeService(ctx, &attributes)
	return attributes
}

func getSellingPriceTaxAttributeService(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	sellingPrice := sellingPriceHandler(ctx)
	if len(sellingPrice) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "list_price_tax",
			Value: "1",
		}
		*attributes = append(*attributes, attribute)
	}
}

func getSpecialKeyAttributeService(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	table, _ := ctx.Page.Locator(".main_inner.item-spec .infoBox.infoBox02.mb20 dl").All()
	for _, item := range table {
		key, _ := item.Locator("dt").InnerText()
		value, _ := item.Locator("dd").InnerText()
		if len(value) > 0 && len(key) > 0 && key != "その他" {
			{
				attribute := ninjacrawler.AttributeItem{
					Key:   key,
					Value: value,
				}
				*attributes = append(*attributes, attribute)
			}
		}
	}
}

func reviewsHandler(ctx ninjacrawler.CrawlerContext) []string {
	var reviews []string
	recentComment := ctx.Document.Find("p.item-review__recent__text").First().Text()
	reviews = append(reviews, recentComment)
	ctx.Document.Find("p.item-review__item__text").Each(func(i int, s *goquery.Selection) {
		reviewText := s.Text()
		reviews = append(reviews, reviewText)
	})
	return reviews
}
