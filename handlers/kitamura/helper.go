package kitamura

import (
	"combined-crawler/pkg/ninjacrawler"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log"
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

func getImagesFromJson(ctx ninjacrawler.CrawlerContext) []string {
	var images []string
	type ProductSchema struct {
		Image []string `json:"image"`
	}

	// Find the script tag containing JSON
	scriptTag := ctx.Document.Find("script[type='application/ld+json']").First()
	jsonContent := scriptTag.Text()

	// Extract the image property using a regular expression
	re := regexp.MustCompile(`"image":\s*\[.*?\]`)
	matches := re.FindString(jsonContent)
	if matches == "" {
		html, _ := ctx.Document.Html()
		ctx.App.Logger.Html(html, ctx.UrlCollection.Url, "image property not found in JSON")
		log.Println("Error: Image property not found in JSON . Url: %s", ctx.UrlCollection.Url)
		return images
	}

	// Create a valid JSON string containing only the image property
	jsonContent = "{" + matches + "}"

	// Parse the JSON content
	var product ProductSchema
	err := json.Unmarshal([]byte(jsonContent), &product)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		return images
	}

	// Return the images from the parsed JSON
	return product.Image
}

func getReviewsService(ctx ninjacrawler.CrawlerContext) []string {
	reviews := []string{}
	reviewListMoreLink, exist := ctx.Document.Find("a.review-list-more-btn").Attr("href")
	if !exist {
		return reviews
	}
	reviewListMoreLink = ctx.App.GetFullUrl(reviewListMoreLink)
	document, err := ctx.App.NavigateToStaticURL(ctx.App.GetHttpClient(), reviewListMoreLink, ctx.App.CurrentProxy)
	if err != nil {
		ctx.App.Logger.Error("reviewListMoreLink: %v", err.Error())
		return reviews
	}
	document.Find(".review-list-area a").Each(func(i int, s *goquery.Selection) {
		reviewLink, ok := s.Attr("href")
		reviewLink = ctx.App.GetFullUrl(reviewLink)
		if ok {
			reviewDocument, reviewErr := ctx.App.NavigateToStaticURL(ctx.App.GetHttpClient(), reviewLink, ctx.App.CurrentProxy)
			if reviewErr != nil {
				ctx.App.Logger.Error("reviewLink: %v", err.Error())
			}
			review := GetDescriptionText(reviewDocument.Find(".review-form-pros-area"))
			reviews = append(reviews, review)
		}
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

	attributes = append(attributes, ninjacrawler.AttributeItem{Key: "selling_price_tax", Value: "1"})
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
