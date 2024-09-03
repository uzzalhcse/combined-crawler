package kitamura

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func getJanService(ctx ninjacrawler.CrawlerContext) string {
	janCode := ""
	re := regexp.MustCompile(`\/pd\/(\d+)`)
	match := re.FindStringSubmatch(ctx.UrlCollection.Url)
	if len(match) > 1 {
		janCode = match[1]
		return janCode
	} else {
		return janCode
	}
}
func getUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}
func getMakerService(ctx ninjacrawler.CrawlerContext) string {
	maker := ""
	makerSelection := ctx.Document.Find("span.secondary-link-text").First()
	maker = strings.TrimSpace(makerSelection.Text())
	return maker
}

func getProductNameService(ctx ninjacrawler.CrawlerContext) string {
	productName := ""
	titleSelection := ctx.Document.Find("h1.product-name")
	if titleSelection.Length() > 0 {
		productName = strings.TrimSpace(titleSelection.Text())
	}

	return productName
}

func getCategoryService(ctx ninjacrawler.CrawlerContext) string {
	catList := []string{}
	categoriesUl := ctx.Document.Find("ul.v-breadcrumbs.breadcrumbs").First()
	if categoriesUl.Length() == 0 {
		return ""
	}

	categoriesUl.Find("li").Each(func(_ int, catItem *goquery.Selection) {
		categoryText := strings.TrimSpace(catItem.Text())
		if categoryText != "" {
			catList = append(catList, categoryText)
		}
	})
	if len(catList) > 2 {
		category := strings.Join(catList, " > ")
		return category
	}
	return ""
}

func getDescriptionService(ctx ninjacrawler.CrawlerContext) string {
	description := ""
	descriptionDiv := ctx.Document.Find("#product-tabs-scroll").Find(".description-area-wide").Find("div").Eq(1)
	// ctx.Document.Find("div[data-v-41f3baee]").First()

	if descriptionDiv.Length() == 0 {
		return ""
	}

	description = GetDescriptionText(descriptionDiv)

	// Remove specific phrases from the description
	description = strings.ReplaceAll(description, "商品説明", "")
	description = strings.ReplaceAll(description, "【製品特徴】", "")
	description = strings.ReplaceAll(description, "※急な欠品や生産完了等でお渡しできない場合はご連絡いたします", "")
	description = strings.ReplaceAll(description, "※アクセサリーの対応機種はこちらのメーカーホームページでご確認ください", "")
	description = strings.ReplaceAll(description, "※アクセサリーの対応機種はココをクリックしてメーカーホームページをご確認ください", "")

	description = strings.TrimSpace(description)

	return description
}

func getReviewsService(ctx ninjacrawler.CrawlerContext) []string {
	reviews := []string{}
	ctx.Document.Find(".review-item-description").Each(func(i int, s *goquery.Selection) {
		review := GetDescriptionText(s)
		reviews = append(reviews, review)
	})
	return reviews
}

func getSellingPriceService(ctx ninjacrawler.CrawlerContext) string {
	sellingPrice := ""
	sellingPriceSelection := ctx.Document.Find("span.sell-price-font.d-inline-block").First()
	sellingPrice = ctx.App.ToNumericsString(sellingPriceSelection.Text())
	return sellingPrice
}

func getAttributeService(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}

	getColorAttributeService(ctx, &attributes)
	getProductDetailsAttributeService(ctx, &attributes)

	return attributes
}

func getColorAttributeService(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	color := ""

	if len(color) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "color",
			Value: color,
		}
		*attributes = append(*attributes, attribute)
	}
}

func getProductDetailsAttributeService(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	// Find the product details element by CSS selector
	productDetailsElement := ctx.Document.Find("#product-tabs-scroll").Find(".description-area-wide").Find("div").Eq(2)

	// Get the text content of the element
	productDetailsText := productDetailsElement.Text()

	// Remove specific phrases from the product details
	productDetailsText = strings.ReplaceAll(productDetailsText, "【製品仕様】", "")
	productDetailsText = strings.ReplaceAll(productDetailsText, "※急な欠品や生産完了等でお渡しできない場合はご連絡いたします", "")
	productDetailsText = strings.ReplaceAll(productDetailsText, "※アクセサリーの対応機種はこちらのメーカーホームページでご確認ください", "")
	productDetailsText = strings.ReplaceAll(productDetailsText, "※アクセサリーの対応機種はココをクリックしてメーカーホームページをご確認ください", "")

	if len(productDetailsText) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "product_details",
			Value: productDetailsText,
		}
		*attributes = append(*attributes, attribute)
	}
}
func GetDescriptionText(selection *goquery.Selection) string {
	textSlice := selection.Contents().Map(func(i int, s *goquery.Selection) string {
		if s.Is("br") {
			return "\n"
		} else {
			text := s.Text()
			text = strings.Trim(text, " \n")
			return text
		}
	})
	text := strings.Join(textSlice, "")
	text = strings.Trim(text, " \n")

	return text
}
