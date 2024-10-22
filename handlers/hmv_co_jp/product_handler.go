package hmv_co_jp

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
				{
					Query: ".singleMainPhotos img",
					Attr:  "src",
				},
			},
		},
		ProductCodes:     []string{},
		Maker:            "",
		Brand:            "",
		ProductName:      &ninjacrawler.SingleSelector{Selector: "h1.title"},
		Category:         "",
		Description:      &ninjacrawler.SingleSelector{Selector: ".singleBasicInfo"},
		Reviews:          []string{},
		ItemTypes:        []string{},
		ItemSizes:        []string{},
		ItemWeights:      []string{},
		SingleItemSize:   "",
		SingleItemWeight: "",
		NumOfItems:       "",
		ListPrice:        "",
		SellingPrice:     &ninjacrawler.SingleSelector{Selector: ".singlePriceBlock .price"},
		Attributes:       GetProductAttribute,
	}
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        productDetailSelector,
			Preference:       ninjacrawler.Preference{ValidationRules: []string{"PageTitle", "SellingPrice", "ProductName"}},
			Engine:           ninjacrawler.Engine{},
		},
	})
}
