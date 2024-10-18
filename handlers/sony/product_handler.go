package sony

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
		Url: GetUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: ".ProductIntroPlate__ThumbImage", Attr: "data-background-image-hires"},
				{Query: ".s5-PDBslideshowA__thumbItem img", Attr: "src"},
				{Query: ".s5-PDBslideshowD__bigGalleryItem-imgDiv img", Attr: "src"},
			},
			IsUnique: true,
		},
		ProductCodes: GetProductCodes,
		Maker:        "",
		Brand:        "",
		ProductName: &ninjacrawler.SingleSelector{
			Selector: "h1.CategoryNav__PdpHeaderTitleName",
		},
		Category: &ninjacrawler.SingleSelector{
			Selector: "h2.CategoryNav__MainName",
		},
		Description: &ninjacrawler.SingleSelector{
			Selector: ".ProductSummary__BodyCopy",
		},
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
			Engine: ninjacrawler.Engine{
				IsDynamic: ninjacrawler.Bool(true),
			},
		},
	})
}
