package kitamura

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
		Url:    getUrlHandler,
		Images: getImagesFromJson,
		ProductCodes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		Maker:       getMakerService,
		Brand:       "",
		ProductName: getProductNameService,
		Category:    getCategoryService,
		Description: getDescriptionService,
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
		ListPrice:        "",
		SellingPrice:     getSellingPriceService,
		Attributes:       getAttributeService,
	}
	crawler.CrawlPageDetail([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        productDetailSelector,
			Preference:       ninjacrawler.Preference{ValidationRules: []string{"Images"}},
			Engine: ninjacrawler.Engine{
				WaitForSelector: ninjacrawler.String("div.product-image-thumbnail-list img"),
				ProviderOption: ninjacrawler.ProviderQueryOption{
					WaitFor: ".v-breadcrumbs__item",
				},
			},
		},
	})
}
