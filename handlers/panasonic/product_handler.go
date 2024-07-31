package panasonic

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan: "",
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: func(ctx ninjacrawler.CrawlerContext) string {
			return ctx.UrlCollection.Url
		},
		Images:       getImagesService,
		ProductCodes: getProductCodesService,
		Maker:        "",
		Brand:        "",
		ProductName:  ProductNameHandler,
		Category:     getCategoryService,
		Description:  getDescriptionService,
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
		SellingPrice:     "",
		Attributes:       GetProductAttribute,
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
