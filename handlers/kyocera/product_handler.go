package kyocera

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {
	crawler.ProductDetailSelector = ninjacrawler.ProductDetailSelector{
		Jan: "",
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: GetUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: ".details .intro .image img", Attr: "src"},
			},
		},
		ProductCodes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		Maker:       "",
		Brand:       "",
		ProductName: ProductNameHandler,
		Category:    GetProductCategory,
		Description: GetProductDescription,
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
	crawler.Collection(constant.ProductDetails).CrawlPageDetail(constant.Products)
}
