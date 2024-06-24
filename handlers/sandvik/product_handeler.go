package sandvik

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan:       JanHandler,
		PageTitle: getPageTitleHandler,
		Url:       getUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: ".img-wrap.position-relative.large-image img", Attr: "src"},
				{Query: ".img-wrap.position-relative.x-large-image img", Attr: "src"},
				{Query: ".row.m-0.mt-4.ng-star-inserted div img", Attr: "src"},
			},
		},
		ProductCodes:     productCodeHandler,
		Maker:            "",
		Brand:            "",
		ProductName:      productNameHandler,
		Category:         getProductCategory,
		Description:      "getProductDescription",
		Reviews:          func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		ItemTypes:        func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		ItemSizes:        func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		ItemWeights:      func(ctx ninjacrawler.CrawlerContext) []string { return []string{} },
		SingleItemSize:   "",
		SingleItemWeight: "",
		NumOfItems:       "",
		ListPrice:        "",
		SellingPrice:     "",
		Attributes:       getProductAttributes,
	}

	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        productDetailSelector,
			Preference:       ninjacrawler.Preference{ValidationRules: []string{"PageTitle"}},
		},
	})
}

func productNameHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.Document.Find(".cor-color-gold-04.my-0.ng-star-inserted").Text()
}

func productCodeHandler(ctx ninjacrawler.CrawlerContext) []string {
	var productCodes []string
	ctx.Document.Find(".cor-font-section-title2.product-title.my-0.ng-star-inserted").Each(func(i int, s *goquery.Selection) {
		productCode := strings.Trim(s.Text(), " \n")
		productCodes = append(productCodes, productCode)
	})
	return productCodes
}

func getPageTitleHandler(ctx ninjacrawler.CrawlerContext) string {
	return strings.Trim(ctx.Document.Find(".cor-font-section-title2.product-title.my-0.ng-star-inserted").First().Text(), " \n")
}

func getUrlHandler(ctx ninjacrawler.CrawlerContext) string {
	return ctx.UrlCollection.Url
}

func getProductCategory(ctx ninjacrawler.CrawlerContext) string {
	categoryItems := make([]string, 0)
	ctx.Document.Find(".breadcrumblist.pt-1.pt-lg-2.pb-4.pb-lg-6 li").Each(func(i int, s *goquery.Selection) {
		if i >= 2 {
			txt := strings.TrimSpace(s.Text())
			if txt == "chevron_right" {
				return
			}
			categoryItems = append(categoryItems, txt)
		}
	})
	return strings.Join(categoryItems, " > ")
}

func JanHandler(ctx ninjacrawler.CrawlerContext) string {
	jan := ctx.Document.Find(".col-12.col-lg-6.px-0.pe-lg-1 product-details-codes-v5 div:nth-child(4)").Text()
	jan = strings.TrimSpace(strings.ReplaceAll(jan, "EAN: ", ""))
	return jan
}

func getProductAttributes(ctx ninjacrawler.CrawlerContext) []ninjacrawler.AttributeItem {
	var attributes []ninjacrawler.AttributeItem
	var attribute ninjacrawler.AttributeItem
	//left Section
	ctx.Document.Find(".col-12.col-lg-6.px-0.pe-lg-1 product-details-codes-v5 div").Each(func(i int, s *goquery.Selection) {

		txt := strings.TrimSpace(s.Text())
		parts := strings.SplitN(txt, ":", 2)
		attribute.Key = strings.TrimSpace(parts[0])
		attribute.Value = strings.TrimSpace(parts[1])
		if attribute.Key != "EAN" {
			attributes = append(attributes, attribute)
		}
	})
	ctx.Document.Find(".product-data.px-1.mb-small.position-relative .row.px-0.border-bottom-gold-02.border-bottom-1.ng-star-inserted").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			key := ""
			s.Find(".col-12.col-md-7.py-2.pb-0.pb-md-2.ps-0 span").Each(func(i int, s *goquery.Selection) {
				key += strings.TrimSpace(s.Text())
			})

			value := s.Find(".col-12.col-md-5.pb-2.ps-0.py-md-2.fw-bold span").First().Text()
			attribute.Key = key
			attribute.Value = value
			if key != "" && value != "" {
				attributes = append(attributes, attribute)
			}
		}
	})
	seeMoreButton := ctx.Page.Locator(".cor-link-button.px-0")
	hasSeeMoreButton, _ := seeMoreButton.IsVisible()
	if hasSeeMoreButton {
		ctx.App.Logger.Info(seeMoreButton.TextContent())
		seeMoreButton.Click()
		time.Sleep(1 * time.Second)
	}
	ctx.Document.Find(".product-data.px-1.mb-small.position-relative .row.px-0.msn-1.border-bottom-gold-01.border-bottom-1.data-transition.ng-star-inserted").Each(func(i int, s *goquery.Selection) {
		key := ""
		s.Find(".col-12.col-md-7.py-2.pb-0.pb-md-2.ps-0 span").Each(func(i int, s *goquery.Selection) {
			key += strings.TrimSpace(s.Text())
		})

		value := s.Find(".col-12.col-md-5.pb-2.ps-0.py-md-2.fw-bold span").First().Text()
		attribute.Key = key
		attribute.Value = value
		if key != "" && value != "" {
			attributes = append(attributes, attribute)
		}
	})

	trialValueSelector := "product-details-start-values-v5 .d-block.mb-small.ng-star-inserted .ng-star-inserted .cor-font-section-title4.pb-3"

	ctx.Document.Find(trialValueSelector).Each(func(i int, s *goquery.Selection) {
		key := s.Text()
		value := ctx.Document.Find("product-details-start-values-v5 .d-block.mb-small.ng-star-inserted div .row.ng-star-inserted").First().Text()
		attribute.Key = key
		attribute.Value = value
		if key != "" && value != "" {
			attributes = append(attributes, attribute)
		}
	})
	return attributes
}
