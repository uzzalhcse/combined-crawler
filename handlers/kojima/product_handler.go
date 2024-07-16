package kojima

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

var ProductDetailSelector ninjacrawler.ProductDetailSelector

func ProductHandler(crawler *ninjacrawler.Crawler) {
	ProductDetailSelector = ninjacrawler.ProductDetailSelector{
		Jan: getJanHandler,
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: getUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: "div.molProductsImages img", Attr: "src"},
			},
		},
		ProductCodes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		Maker: "",
		Brand: "",
		ProductName: &ninjacrawler.SingleSelector{
			Selector: "h1.name.opt-large.just-bold.mt2",
		},
		Category:    getProductCategory,
		Description: getProductDescription,
		Reviews:     getReviewsService,
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
		SellingPrice:     getSellingPriceService,
		Attributes:       getProductAttribute,
	}
	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        ProductDetailSelector,
			Preference:       ninjacrawler.Preference{ValidationRules: []string{"PageTitle"}},
		},
	})
}
