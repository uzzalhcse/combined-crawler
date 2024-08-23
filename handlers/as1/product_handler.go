package as1

import (
	"combined-crawler/constant"
	"combined-crawler/pkg/ninjacrawler"
)

func ProductHandler(crawler *ninjacrawler.Crawler) {
	productDetailSelector := ninjacrawler.ProductDetailSelector{
		Jan: getJanService,
		PageTitle: &ninjacrawler.SingleSelector{
			Selector: "title",
		},
		Url: getUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: "#gallery > ul.gallery-thumbnails > li img", Attr: "src"},
				{Query: "#gallery > ul > li > a > img", Attr: "src"},
			},
		},
		ProductCodes: getProductCode,
		Maker:        getMaker,
		Brand:        "",
		ProductName:  getProductNameService,
		Category:     "",
		Description:  "",
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
		SellingPrice:     getSellingPriceService,
		Attributes:       getProductAttribute,
	}
	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        productDetailSelector,
			Preference:       ninjacrawler.Preference{ValidationRules: []string{"PageTitle"}},
			Engine: ninjacrawler.Engine{
				ProviderOption: ninjacrawler.ProviderQueryOption{
					Wait: 10000,
				},
			},
		},
	})
}
