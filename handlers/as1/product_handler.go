package as1

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
		Url: getUrlHandler,
		Images: &ninjacrawler.MultiSelectors{
			Selectors: []ninjacrawler.Selector{
				{Query: ".goods_detail a.thumbnail img", Attr: "src"},
				{Query: ".goods_c img.img-responsive", Attr: "src"},
				{Query: ".goods_a img.img-responsive", Attr: "src"},
			},
		},
		ProductCodes: getProductCode,
		Maker:        getMaker,
		Brand:        "",
		ProductName: &ninjacrawler.SingleSelector{
			Selector: ".goodsname01",
			Regexp:   []string{`\s+`},
		},
		Category:    productCategoryHandler,
		Description: getProductDescription,
		Reviews: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		ItemTypes: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		ItemSizes: getItemSizes,
		ItemWeights: func(ctx ninjacrawler.CrawlerContext) []string {
			return []string{}
		},
		SingleItemSize:   "",
		SingleItemWeight: "",
		NumOfItems:       "",
		ListPrice:        "",
		SellingPrice: &ninjacrawler.SingleSelector{
			Selector: "td span.goodsprice",
			Regexp:   []string{`[^0-9]`},
		},
		Attributes: getProductAttribute,
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
