package rakuten

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
					Query: ".oneImageWrap img",
					Attr:  "src",
				},
			},
		},
		ProductCodes:     []string{},
		Maker:            "Panasonic",
		Brand:            "",
		ProductName:      &ninjacrawler.SingleSelector{Selector: "#productTitle"},
		Category:         "",
		Description:      &ninjacrawler.SingleSelector{Selector: ".linkOtherFormat ul"},
		Reviews:          []string{},
		ItemTypes:        []string{},
		ItemSizes:        []string{},
		ItemWeights:      []string{},
		SingleItemSize:   "",
		SingleItemWeight: "",
		NumOfItems:       "",
		ListPrice:        "",
		SellingPrice:     &ninjacrawler.SingleSelector{Selector: ".productPrice .price"},
		Attributes:       GetProductAttribute,
	}
	crawler.Crawl([]ninjacrawler.ProcessorConfig{
		{
			Entity:           constant.ProductDetails,
			OriginCollection: constant.Products,
			Processor:        productDetailSelector,
			Preference:       ninjacrawler.Preference{ValidationRules: []string{"PageTitle|required|blacklists:楽天ブックス: お探しのページが見つかりません"}},
			Engine: ninjacrawler.Engine{
				SimulateMouse: ninjacrawler.Bool(false),
			},
		},
	})
}
