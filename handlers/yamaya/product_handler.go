package yamaya

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
	"strings"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan: func(ctx ninjacrawler.CrawlerContext) string {
			index := strings.LastIndex(ctx.UrlCollection.Url, "=")
			janCode := ctx.UrlCollection.Url[index+1:]
			return janCode

		},
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: getUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: ".card-img.mx-auto.d-block", Attr: "src"}, //need to be unique
			},
			ExcludeString: []string{"noimage_M.jpg"},
		},
		ProductCodes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		Maker:       "",
		Brand:       "",
		ProductName: productNameHandler,
		Category:    getProductCategory,
		Description: getProductDescription,
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
		ListPrice:        "",
		SellingPrice:     getSellingPrice,
		Attributes:       getProductAttribute,
	}
	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        productDetailSelector,
			Preference: ninjacrawler.Preference{
				ValidationRules: []string{"PageTitle|blacklists:403 Forbidden,502 Gateway Timeout|required", "ProductName|required"},
			},
		},
	})
}
