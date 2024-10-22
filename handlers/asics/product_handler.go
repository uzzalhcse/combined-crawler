package asics

import (
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func ProductDetailsHandler(crawler *ninjacrawler.Crawler) {
	ProductDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan:       "",
		PageTitle: "pageTitleHandler",
		Url: func(ctx ninjacrawler.CrawlerContext) string {
			return ctx.UrlCollection.Url
		},
		Images: &ninjacrawler.MultiSelectors{
			//Selectors: []ninjacrawler.Selector{
			//	//{Query: "div.swiper-slide.u-aspect.is-1-1.u-mb4 img", Attr: "src"},
			//},
		},
		ProductCodes: productCodeHandler,
		Maker: func(ctx ninjacrawler.CrawlerContext) string {
			return "Asics"
		},
		Brand: func(ctx ninjacrawler.CrawlerContext) string {
			return "Asics"
		},
		ProductName:      "productNameHandler",
		Category:         "productCategoryHandler",
		Description:      "productDescriptionHandler",
		Reviews:          func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		ItemTypes:        func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		ItemSizes:        func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		ItemWeights:      func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		SingleItemSize:   "",
		SingleItemWeight: "",
		NumOfItems:       "",
		ListPrice:        "",
		SellingPrice:     "sellingPriceHandler",
		Attributes:       productAttributesHandler,
	}

	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           ProductDetails,
			OriginCollection: Products,
			Processor:        ProductDetailSelector,
		},
	})
}

func pageTitleHandler(ctx ninjacrawler.CrawlerContext) string {
	pageTitle := ctx.Document.Find("title").Text()
	return pageTitle
}

func productNameHandler(ctx ninjacrawler.CrawlerContext) string {
	title := ctx.Document.Find("h1.comp-title.u-heading-l.u-fw-black.u-mt16.u-mb16.u-size20-sp").Text()
	return title
}

func productCategoryHandler(ctx ninjacrawler.CrawlerContext) string {
	category := ""
	ctx.Document.Find("ol.u-breadcrumb-txt-pc.u-breadcrumb-txt-sp li a span").Each(func(_ int, breadcrumb *goquery.Selection) {
		if category != "" {
			category += " > "
		}
		category += breadcrumb.Text()
	})
	return category
}

func productCodeHandler(ctx ninjacrawler.CrawlerContext) []string {
	return []string{""}
}

func productDescriptionHandler(ctx ninjacrawler.CrawlerContext) string {
	description := ctx.Document.Find("#js-oveflowAdjustBody p").Text()
	description = strings.ReplaceAll(description, " ", "")
	return description
}

func sellingPriceHandler(ctx ninjacrawler.CrawlerContext) string {
	sellingPrice := ctx.Document.Find("span.comp-price.u-size32.u-lh-tight.u-fontEn.u-size24-sp").Text()
	return ctx.App.ToNumericsString(sellingPrice)
}

func productAttributesHandler(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	var attributes []ninjacrawler.AttributeItem
	attributes = append(attributes, ninjacrawler.AttributeItem{
		Key:   "selling_price_tax",
		Value: "1",
	})
	//ctx.Page.WaitForSelector(".swiper-slide.u-h-auto.u-relative.u-w-88-pc.swiper-slide-active")
	//colorSwitchButtons, _ := ctx.Page.Locator("div.swiper-slide.u-h-auto.u-relative.u-w-88-pc").All()
	//
	//for ind, _ := range colorSwitchButtons {
	//	if ind != 0 {
	//		err := ctx.Page.Locator("div.swiper-slide.u-h-auto.u-relative.u-w-88-pc:nth-child(" + strconv.Itoa(ind+1) + ")").Click()
	//		time.Sleep(time.Millisecond * 100)
	//		if err != nil {
	//			ctx.App.Logger.Error("can't click on color switch button", err)
	//		}
	//	}
	//	text, _ := ctx.Page.Locator("div.comp-color-ttl.u-size16.u-fw-black.u-lh-normal.u-mb8").InnerText()
	//	attributes = append(attributes, ninjacrawler.AttributeItem{
	//		Key:   "color_variations",
	//		Value: text,
	//	})
	//}
	return attributes
}
