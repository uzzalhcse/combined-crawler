package kojima

import (
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func getJanHandler(ctx ninjacrawler.CrawlerContext) string {
	janCode := ""
	re := regexp.MustCompile(`\bprod=(\d+)`)
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

func getProductCategory(ctx ninjacrawler.CrawlerContext) string {
	catList := []string{}
	ctx.Document.Find("div.Breadcrumb.MK2PFRDH300_01 li").Each(func(_ int, catItem *goquery.Selection) {
		categoryText := strings.TrimSpace(catItem.Text())
		if categoryText != "" {
			catList = append(catList, categoryText)
		}
	})
	if len(catList) > 2 {
		category := strings.Join(catList[1:len(catList)-1], " > ")
		return category
	}
	return ""
}

func getProductDescription(ctx ninjacrawler.CrawlerContext) string {
	description := ctx.Document.Find("div.ProductClassTable.MK2PFRPM000_02 div.sub-text").Text()
	return strings.TrimSpace(description)
}

func getReviewsService(ctx ninjacrawler.CrawlerContext) []string {
	reviews := []string{}

	reviewHyperlink, exists := ctx.Document.Find("a.molButton.plain.shadow.h-small.icon-arrow-black").Attr("href")
	if !exists {
		return reviews
	}
	reviewUrl := ctx.App.GetFullUrl(reviewHyperlink)
	for {
		httpClient := ctx.App.GetHttpClient()
		pageData, err := ctx.App.NavigateToStaticURL(httpClient, reviewUrl, ninjacrawler.Proxy{})
		if err != nil {
			ctx.App.Logger.Warn(err.Error())
			return reviews

		}
		reviewItems := pageData.Find("p.normal-text.mt3")
		reviewItems.Each(func(_ int, reviewItem *goquery.Selection) {
			reviewText := strings.TrimSpace(reviewItem.Text())
			reviews = append(reviews, reviewText)
		})

		nextPage := pageData.Find("a.next")
		if nextPage.Length() == 0 {
			break
		}
		nextPageURL, exist := nextPage.Attr("href")
		if !exist {
			fmt.Println("Next page URL not found.")
			break
		}
		reviewUrl = ctx.App.GetFullUrl(nextPageURL)
	}
	return reviews
}

func getListPriceService(ctx ninjacrawler.CrawlerContext) string {
	listPriceValue := ""
	keyFound := false
	tables := ctx.Document.Find("table.data-table.molTableSimple")
	tables.Each(func(_ int, table *goquery.Selection) {
		thData := []string{}
		tdData := []string{}

		thElements := table.Find("tbody th")
		thElements.Each(func(_ int, th *goquery.Selection) {
			thData = append(thData, strings.TrimSpace(th.Text()))
		})
		tdElements := table.Find("td")
		tdElements.Each(func(_ int, td *goquery.Selection) {
			tdData = append(tdData, strings.TrimSpace(td.Text()))
		})

		//initialize a blank map/dictionary
		spec_dic := make(map[string]string)

		for i, key := range thData {
			value := tdData[i]
			spec_dic[key] = value
		}
		if val, ok := spec_dic["メーカー希望小売価格"]; ok {
			listPriceValue = val
			keyFound = true
		}
	})

	if keyFound {
		return listPriceValue
	}

	return ""
}

func getSellingPriceService(ctx ninjacrawler.CrawlerContext) string {
	priceSection := ctx.Document.Find("dl.molProductDefinition.w180.pt2.mt3").First()
	priceSection = priceSection.Find("span.number").First()
	priceText := priceSection.Text()
	price := ctx.App.ToNumericsString(priceText)

	return price

}
func getProductAttribute(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	attributes := []ninjacrawler.AttributeItem{}
	getSellingPriceTaxAttributeService(ctx, &attributes)
	getSpecData(ctx, &attributes)
	return attributes
}

func getSellingPriceTaxAttributeService(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	sellingPrice := getSellingPriceService(ctx)
	if len(sellingPrice) > 0 {
		attribute := ninjacrawler.AttributeItem{
			Key:   "selling_price_tax",
			Value: "1",
		}
		*attributes = append(*attributes, attribute)
	}
}

func getSpecData(ctx ninjacrawler.CrawlerContext, attributes *[]ninjacrawler.AttributeItem) {
	tables := ctx.Document.Find("table.data-table.molTableSimple")
	tables.Each(func(_ int, table *goquery.Selection) {
		thData := []string{}
		tdData := []string{}

		thElements := table.Find("tbody th")
		thElements.Each(func(_ int, th *goquery.Selection) {
			thData = append(thData, strings.TrimSpace(th.Text()))
		})
		tdElements := table.Find("td")
		tdElements.Each(func(_ int, td *goquery.Selection) {
			tdData = append(tdData, strings.TrimSpace(td.Text()))
		})

		if len(thData) > 0 && len(thData) == len(tdData) {
			for i := range thData {
				*attributes = append(*attributes, ninjacrawler.AttributeItem{
					Key:   thData[i],
					Value: tdData[i],
				})
			}
		}
	})
}
