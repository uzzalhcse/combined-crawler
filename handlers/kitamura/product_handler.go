package kitamura

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {

	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        handleProduct,
			Preference: ninjacrawler.Preference{
				ValidationRules: []string{"PageTitle|required"},
			},
		},
	})
}
func handleProduct(ctx ninjacrawler.CrawlerContext, fn func([]ninjacrawler.ProductDetailSelector, string)) error {
	urlCollections := []ninjacrawler.UrlCollection{}
	ctx.Document.Find("dl#product_detail_special span.product_detail_container a").Each(func(i int, a *goquery.Selection) {
		href, exists := a.Attr("href")
		if exists {
			href = ctx.App.GetFullUrl(href)
			if !isValidHost(href) {
				fmt.Println("Invalid host: ", href)
				return
			}
			urlCollections = append(urlCollections, ninjacrawler.UrlCollection{
				Url:      href,
				Parent:   ctx.UrlCollection.Url,
				MetaData: ctx.UrlCollection.MetaData,
			})
		}
	})

	if len(urlCollections) > 0 {
		ctx.App.InsertUrlCollections(constant.Products, urlCollections, ctx.UrlCollection.Url)
	}
	fn([]ninjacrawler.ProductDetailSelector{
		selectProduct(),
	}, "")
	return nil

}

func selectProduct() ninjacrawler.ProductDetailSelector {
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan: getJanService,
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: GetUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: "p#product_img img", Attr: "src"},
				{Query: "p.product_img img", Attr: "src"},
			},
		},
		ProductCodes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		Maker: func(ctx ninjacrawler.CrawlerContext) string {
			return "Suntory"

		},
		Brand:       "",
		ProductName: getProductNameService,
		Category:    getCategoryService,
		Description: getDescriptionService,
		Reviews: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		ItemTypes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		ItemSizes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		ItemWeights: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		SingleItemSize:   "",
		SingleItemWeight: "",
		NumOfItems:       "",
		ListPrice:        getListPriceService,
		SellingPrice:     "",
		Attributes:       getProductAttribute,
	}
	return productDetailSelector
}
